package schema

import (
	"fmt"
	"go/ast"
	"idiocy/apitmpl"
	"idiocy/helper"
	"idiocy/logger"
	"idiocy/platform"
	"strings"

	"github.com/fatih/structtag"
	"github.com/getkin/kin-openapi/openapi3"
)

func (f *SourceFile) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.AstFile)
}

func (f *SourceFile) find() {
	f.BuildStacks()
	f.FindGinInstance()
}

func (f *SourceFile) FindGinInstance() {
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}
		identExpr, identOK := node.(*ast.Ident)
		callExpr, callOK := node.(*ast.CallExpr) // CallExpr
		switch {
		case identOK:
			if identExpr.Obj != nil {
				if identExpr.Obj.Decl != nil {
				}
			}

		case callOK:
			if !DetectInstanceCall(callExpr.Fun, platform.GinInstanceCall) {
				return false
			}
			index := f.NodeIndex(node)
			f.FindCallLIdent(index)

		default:
			logger.S.Infof(logger.ColorLightGreen("unhandle %d --> %d"), node.Pos(), node.End())
			logger.S.Info(logger.ColorLightGreen(string(f.Content[node.Pos()-1 : node.End()-1])))
		}

		return true
	})
}

// helpers
// =======

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

// exports
// =======

func (f *SourceFile) EnumerateGinHandles() {
	if len(f.GinIdents) == 0 {
		return
	}
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		callExpr, callExprOK := n.(*ast.CallExpr)
		switch {
		case callExprOK:
			if selectorExpr, selectorExprOK := callExpr.Fun.(*ast.SelectorExpr); selectorExprOK {
				_ = selectorExpr
				if x, sel, ok := extractSelectorExpr(callExpr.Fun.(*ast.SelectorExpr)); ok {
					_ = x
					switch sel.Name {
					case "GET", "PUT", "PATH", "POST", "DELETE":
					}
				}

			}
		}
		return true
	})
}

func (f *SourceFile) EnumerateGinBind(Lbrace, Rbrace int) (paramName string) {
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		from, to := 0, 0
		if n.Pos().IsValid() {
			from = int(n.Pos())
			if from < Lbrace {
				return true
			}
		}

		if n.End().IsValid() {
			to = int(n.End())
			if to >= Rbrace {
				return true
			}
		}

		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex
		fun, args, ok := helper.ExtractCallExpr(n)
		_ = args
		if !ok {
			return true
		}

		x, sel, ok := helper.ExtractSelectorExpr(ast.Node(*fun))
		_ = x
		if !ok {
			return true
		}

		switch sel.Name {
		case "ShouldBind":
			nn := f.fullStacks[nodeIndex+6]
			if _, obj, ok := helper.ExtractIdent(nn); ok {
				fullname := helper.ExtractObjectTypeName(obj)
				paramName = fullname
				return false
			}
		}

		return true
	})
	return
}

func (f *SourceFile) EnumerateGinRoute(groupRoute string, Lbrace, Rbrace int) {
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		from, to := 0, 0
		if n.Pos().IsValid() {
			from = int(n.Pos())
			if from < Lbrace {
				return true
			}
		}

		if n.End().IsValid() {
			to = int(n.End())
			if to >= Rbrace {
				return true
			}
		}

		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex

		if nextSel, ok := n.(*ast.SelectorExpr); ok {
			nx, nsel, _ := extractSelectorExpr(nextSel)
			_ = nx
			switch nsel.Name {
			case "GET", "POST", "PUT", "PATCH", "DELETE":
				routePathNode := f.fullStacks[nodeIndex+3]
				routePath := ""
				if basicLit, basicLitOK := routePathNode.(*ast.BasicLit); basicLitOK {
					if strings.EqualFold(groupRoute, "") {
						routePath = basicLit.Value
					} else {
						routePath = fmt.Sprintf(`"%s%s"`, strings.ReplaceAll(groupRoute, `"`, ``), strings.ReplaceAll(basicLit.Value, `"`, ``))
					}
				}
				if strings.Contains(routePath, "swagger") {
					return true
				}
				route := GinAPI{
					Source:   f,
					Group:    "",
					Method:   nsel.Name,
					Path:     routePath,
					Node:     &n,
					APIParam: &APIParam{},
				}
				ProjSchema.GinAPIs = append(ProjSchema.GinAPIs, &route)
				handlePathNode := f.fullStacks[nodeIndex+4]
				if funcLit, funcLitOk := handlePathNode.(*ast.FuncLit); funcLitOk {
					body := funcLit.Body
					codes, respNames := f.EnumerateGinResponse(int(body.Lbrace), int(body.Rbrace))
					route.APIParam.StructName = f.EnumerateGinBind(int(body.Lbrace), int(body.Rbrace))
					for i := 0; i < len(codes); i++ {
						route.APIResopnse = append(route.APIResopnse, &APIResopnse{
							Code:       codes[i],
							StructName: respNames[i],
						})
					}
					for _, v := range body.List {
						if declStmt, declStmtOK := v.(*ast.DeclStmt); declStmtOK {
							_ = declStmt
						}
					}
				} else {
				}
			}
		}

		return true
	})
	return
}

func (f *SourceFile) EnumerateGinResponse(Lbrace, Rbrace int) (codes, structNames []string) {
	codes, structNames = []string{}, []string{}
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		from, to := 0, 0
		if n.Pos().IsValid() {
			from = int(n.Pos())
			if from < Lbrace {
				return true
			}
		}

		if n.End().IsValid() {
			to = int(n.End())
			if to >= Rbrace {
				return true
			}
		}

		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex
		fun, args, ok := helper.ExtractCallExpr(n)
		if !ok {
			return true
		}

		x, sel, ok := helper.ExtractSelectorExpr(ast.Node(*fun))
		_ = x
		if !ok {
			return true
		}

		switch sel.Name {
		case "JSON":
			statusCode := ""
			responseTypeName := ""
			if len(args) >= 1 {
				if _, value, ok := helper.ExtractBasicLit(ast.Node(args[0])); ok {
					statusCode = value
				}
			}
			if len(args) >= 2 {
				if _, object, ok := helper.ExtractIdent(ast.Node(args[1])); ok {
					responseTypeName = helper.ExtractObjectTypeName(object)
				}
			}
			codes = append(codes, statusCode)
			structNames = append(structNames, responseTypeName)
		}

		return true
	})
	return
}

func (f *SourceFile) EnumerateLazy(prefix string, Lbrace, Rbrace int) (paths, methods, codes, structNames []string) {
	paths, methods, codes, structNames = []string{}, []string{}, []string{}, []string{}
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		from, to := 0, 0
		if n.Pos().IsValid() {
			from = int(n.Pos())
			if from < Lbrace {
				return true
			}
		}

		if n.End().IsValid() {
			to = int(n.End())
			if to >= Rbrace {
				return true
			}
		}

		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex

		if fl, flOK := helper.ExtractFuncLit(n); flOK {
			pathNode := f.fullStacks[nodeIndex-1]
			methodNode := f.fullStacks[nodeIndex-2]
			methodName, _, _ := helper.ExtractIdent(methodNode)

			_, pathString, _ := helper.ExtractBasicLit(pathNode)
			paths = append(paths, strings.ReplaceAll(fmt.Sprintf("%s%s", prefix, pathString), `"`, ``))
			methods = append(methods, methodName)
			codes = append(codes, "200")

			if bl, ok := helper.ExtractBlockStml(fl.Body); ok {
				modelName := f.EnumerateLazyModel(int(bl.Lbrace), int(bl.Rbrace))
				structNames = append(structNames, modelName)
			}
		}

		return true
	})

	return
}

func (f *SourceFile) EnumerateLazyModel(Lbrace, Rbrace int) (modelname string) {
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		from, to := 0, 0
		if n.Pos().IsValid() {
			from = int(n.Pos())
			if from < Lbrace {
				return true
			}
		}

		if n.End().IsValid() {
			to = int(n.End())
			if to >= Rbrace {
				return true
			}
		}

		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex

		//Model
		if name, obj, ok := helper.ExtractIdent(n); ok && strings.Compare(name, "Model") == 0 {
			_ = obj
			nn := f.fullStacks[nodeIndex+1]
			if unary, ok := helper.ExtractUnary(nn); ok {
				if unary.Op == 17 {
					if cl, clOK := helper.ExtractCompositeLit(unary.X); clOK {
						if x, sel, ok := helper.ExtractSelectorExpr(cl.Type); ok {
							name, _, _ := helper.ExtractIdent(ast.Node(*x))
							modelname = fmt.Sprintf("%s.%s", name, sel.Name)
							return false
						}

					}

				}
			}
		}
		return true
	})
	return modelname
}

func (f *SourceFile) EnumerateStructAndGinVars() {
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex
		structType, structTypeOK := n.(*ast.StructType)
		typeSpec, typeSpecOK := n.(*ast.TypeSpec)
		callExpr, callExprOK := n.(*ast.CallExpr)
		ident, identOK := n.(*ast.Ident)
		decl, declOK := n.(*ast.DeclStmt)
		blk, blkOK := n.(*ast.BlockStmt)
		_ = decl
		_ = blk

		switch {
		case blkOK:
			nn := f.fullStacks[nodeIndex-18]

			groupPath := ""

			if name, _, ok := helper.ExtractIdent(nn); ok && strings.EqualFold(name, "Group") {
				nameNode := f.fullStacks[nodeIndex-17]
				if kind, value, ok := helper.ExtractBasicLit(nameNode); ok && kind == 9 {
					groupPath = value
					paths, methods, codes, models := f.EnumerateLazy(groupPath, int(blk.Lbrace), int(blk.Rbrace))
					for i := 0; i < len(paths); i++ {
						switch methods[i] {
						case "GET", "POST", "PUT", "PATCH", "DELETE":
							routePath := paths[i]
							route := GinAPI{
								Source:   f,
								Group:    "",
								Method:   methods[i],
								Path:     routePath,
								Node:     &n,
								APIParam: &APIParam{},
							}
							ProjSchema.GinAPIs = append(ProjSchema.GinAPIs, &route)
							route.APIParam.StructName = models[i]
							route.APIResopnse = append(route.APIResopnse, &APIResopnse{
								Code:       codes[i],
								StructName: models[i],
							})
						}
					}
				}
			}

		case declOK:
		case identOK:
			if ident.Obj != nil {
				if ident.Obj.Kind == ast.Var {
					if assignStmt, assignStmtOK := ident.Obj.Decl.(*ast.AssignStmt); assignStmtOK {
						if len(assignStmt.Rhs) != 0 {
							if callExpr, callExprOK := assignStmt.Rhs[0].(*ast.CallExpr); callExprOK {
								if selectorExpr, selectorExprOK := callExpr.Fun.(*ast.SelectorExpr); selectorExprOK {
									if selectorExpr.Sel.Name == "Group" {
										groupBlockNode := f.fullStacks[nodeIndex+5]
										if kind, value, ok := helper.ExtractBasicLit(groupBlockNode); ok && kind == 9 {
											groupPath := value
											if blk, blkOK := helper.ExtractBlockStml(f.fullStacks[nodeIndex+6]); blkOK {
												f.EnumerateGinRoute(groupPath, int(blk.Lbrace), int(blk.Rbrace))
											}
											newGinIdent := NewGinIdentifier()
											newGinIdent.Source = f
											newGinIdent.Node = ident
											newGinIdent.InstancingCall = callExpr
											existGinIndent := ProjSchema.GinIdentifierWithFileIdent(f, callExpr)
											if existGinIndent == nil {
												ProjSchema.AddGinIdentifier(f, newGinIdent)
												existGinIndent = ProjSchema.GinIdentifierWithFileIdent(f, callExpr)
											}
										}
									}
									if equalSelectorExpr(selectorExpr, "gin", "Default") {
										newGinIdent := NewGinIdentifier()
										newGinIdent.Source = f
										newGinIdent.Node = ident
										newGinIdent.InstancingCall = callExpr
										existGinIndent := ProjSchema.GinIdentifierWithFileIdent(f, callExpr)
										if existGinIndent == nil {
											ProjSchema.AddGinIdentifier(f, newGinIdent)
											existGinIndent = ProjSchema.GinIdentifierWithFileIdent(f, callExpr)
										}
										existGinIndent.AddCall(callExpr)
										callSel := f.fullStacks[nodeIndex-1]
										if nextSel, ok := callSel.(*ast.SelectorExpr); ok {
											nx, nsel, _ := extractSelectorExpr(nextSel)
											_ = nx
											switch nsel.Name {
											case "GET", "POST", "PUT", "PATCH", "DELETE":
												routePathNode := f.fullStacks[nodeIndex+2]
												routePath := ""
												if basicLit, basicLitOK := routePathNode.(*ast.BasicLit); basicLitOK {
													routePath = basicLit.Value
												}
												if strings.Contains(routePath, "swagger") {
													return true
												}
												route := GinAPI{
													Source:   f,
													Group:    "",
													Method:   nsel.Name,
													Path:     routePath,
													Node:     &n,
													APIParam: &APIParam{},
												}
												ProjSchema.GinAPIs = append(ProjSchema.GinAPIs, &route)
												handlePathNode := f.fullStacks[nodeIndex+3]
												if funcLit, funcLitOk := handlePathNode.(*ast.FuncLit); funcLitOk {
													body := funcLit.Body
													codes, respNames := f.EnumerateGinResponse(int(body.Lbrace), int(body.Rbrace))
													route.APIParam.StructName = f.EnumerateGinBind(int(body.Lbrace), int(body.Rbrace))
													for i := 0; i < len(codes); i++ {
														route.APIResopnse = append(route.APIResopnse, &APIResopnse{
															Code:       codes[i],
															StructName: respNames[i],
														})
													}
													for _, v := range body.List {
														if declStmt, declStmtOK := v.(*ast.DeclStmt); declStmtOK {
															_ = declStmt
														}
													}
												} else {
												}

											}
										}

									}
								}
							}
						}
					}

				}
			}

		case callExprOK:
			_ = callExpr
		case typeSpecOK:
			typePos := -1
			if typeSpec.Type.Pos().IsValid() {
				typePos = int(typeSpec.Type.Pos())
			}
			Structs = append(Structs, Struct{f.FullPath, typeSpec.Name.Name, typePos, int(typeSpec.Name.NamePos), f.AstFile, &n})

		case structTypeOK:
			_ = structType
			pNode := f.fullStacks[nodeIndex-1]

			var name string

			name, _, _ = helper.ExtractIdent(pNode)

			doc := apitmpl.Doc
			def := openapi3.SchemaRef{}
			def.Value = openapi3.NewArraySchema()
			def.Value.Type = "object"
			def.Value.Properties = openapi3.Schemas{}

			for _, v := range structType.Fields.List {
				typeDesc, _, _ := helper.ExtractIdent(v.Type)

				switch typeDesc {
				case "int", "uint":
					typeDesc = "integer"
				case "string", "bool":
				default:
					typeDesc = "string"
				}

				propName := ""
				if len(v.Names) > 0 {
					propName = v.Names[0].Name
				}

				if v.Tag != nil {
					vv := strings.TrimPrefix(v.Tag.Value, "`")
					vv = strings.TrimSuffix(vv, "`")
					if tags, err := structtag.Parse(vv); err == nil {
						if jsonTag, err := tags.Get("json"); err == nil {
							_ = jsonTag
							if strings.EqualFold(jsonTag.Name, "-") {
								propName = ""
							} else {
								propName = jsonTag.Name
							}
						} else {
						}
					}
				}
				if !strings.EqualFold(propName, "") {
					ref := openapi3.NewStringSchema().NewRef()
					ref.Value.Type = typeDesc
					ref.Value.Description = fmt.Sprintf("TODO: 缺少 %s 的描述", propName)
					def.Value.Properties[propName] = ref
				}

			}

			if !strings.EqualFold(name, "") {
				doc.Definitions[name] = &def
			}

			sentry := false
			for _, v := range Structs {
				if v.FullPath == f.FullPath && v.TypePos == int(structType.Pos()) {
					sentry = true
					break
				}
			}

			if !sentry {
				Structs = append(Structs, Struct{f.FullPath, "", -1, int(structType.Pos()), f.AstFile, &n})
			}
		}

		return true
	})

}

func (f *SourceFile) BuildStacks() {
	f.fullStacks = []ast.Node{}
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}
		f.fullStacks = append(f.fullStacks, node)
		return true
	})
}

func (f *SourceFile) NodeIndex(node ast.Node) int {
	for k, v := range f.fullStacks {
		switch {
		case v.Pos().IsValid() != node.Pos().IsValid():
			fallthrough
		case v.End().IsValid() != node.End().IsValid():
			continue
		}

		if node.Pos().IsValid() {
			if int(v.Pos()) != int(node.Pos()) {
				continue
			}
		}
		if node.End().IsValid() {
			if int(v.End()) != int(node.End()) {
				continue
			}
		}
		return k
	}
	return -1
}

func (f *SourceFile) StacksLength() int {
	return len(f.fullStacks)
}

func (f *SourceFile) PrintNode(callIndex int) {
	node := f.fullStacks[callIndex]
	logger.S.Debug(string(f.Content[node.Pos()-1 : node.End()-1]))
}

// /ast.Ident
func (f *SourceFile) FindCallLIdent(callIndex int) ast.Node {
	lIndex := callIndex - 1
	llIndex := callIndex - 2
	if lIndex < 0 || llIndex < 0 {
		return nil
	}

	lNode := f.fullStacks[lIndex]
	llNode := f.fullStacks[llIndex]

	_, identOK := lNode.(*ast.Ident)                 // ObjKind.var
	assignStmt, assignOK := llNode.(*ast.AssignStmt) // Lhs,IsOperator
	if !assignOK || !identOK {
		return nil
	}

	if !assignStmt.Tok.IsOperator() {
		return nil
	}

	switch assignStmt.Tok.String() {
	case "=", ":=":
		return lNode
	}

	return nil
}

func DetectInstanceCall(expr ast.Expr, calls []platform.InstanceCall) bool {
	for _, v := range calls {
		if isFuncCallWithName(expr, v.Pkg, v.FuncName) {
			return true
		}
	}
	return false
}

func isFuncCallWithName(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}
