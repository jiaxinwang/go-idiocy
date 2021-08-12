package schema

import (
	"go/ast"
	"go/token"
	"idiocy/logger"
	"idiocy/platform"
	"log"
)

func (f *SourceFile) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.AstFile)
}

func (f *SourceFile) find() {
	f.BuildStacks()
	f.FindDecals()
	f.FindGinInstance()
}

func (f *SourceFile) FindDecals() {
	for _, d := range f.AstFile.Decls {
		if genDecl, ok := d.(*ast.GenDecl); ok {
			if genDecl.Tok == token.CONST || genDecl.Tok == token.VAR {
				for _, cDecl := range genDecl.Specs {
					if vSpec, ok := cDecl.(*ast.ValueSpec); ok {
						for i := 0; i < len(vSpec.Names); i++ {
							// TODO: only basic literals work currently
							switch v := vSpec.Values[i].(type) {
							case *ast.BasicLit:
								f.Decls[vSpec.Names[i].Name] = v.Value
							default:
								logger.S.Infof("Name: %s - Unsupported ValueSpec: %+v\n", vSpec.Names[i].Name, v)
							}
						}
					}
				}
			}

		}
	}

	log.Println("Decls:", f.Decls)
}

func (f *SourceFile) FindGinInstance() {

	// var lastAssignStmt *ast.AssignStmt
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}

		logger.S.Infof("----- %d --> %d -----", int(node.Pos()), int(node.End()))
		logger.S.Info(string(f.Content[node.Pos()-1 : node.End()-1]))

		// BadExpr
		// Ident
		identExpr, identOK := node.(*ast.Ident)
		// Ellipsis
		// BasicLit

		//!FuncLit A function literal represents an anonymous function.
		funcLitExpr, funcLitOK := node.(*ast.FuncLit) // FuncLit
		// CompositeLit
		// ParenExpr
		// SelectorExpr
		// IndexExpr
		// SliceExpr
		// TypeAssertExpr

		callExpr, callOK := node.(*ast.CallExpr) // CallExpr

		// StarExpr
		// UnaryExpr
		// BinaryExpr
		// KeyValueExpr

		// BadStmt
		// DeclStmt
		// EmptyStmt
		// LabeledStmt
		// ExprStmt
		// SendStmt
		// IncDecStmt

		assignStmt, assignOK := node.(*ast.AssignStmt) // AssignStmt
		// GoStmt
		// DeferStmt
		// ReturnStmt
		// BranchStmt

		blockStmt, blockStmtOK := node.(*ast.BlockStmt) // BlockStmt
		// IfStmt
		// CaseClause
		// SwitchStmt
		// TypeSwitchStmt
		// CommClause
		// SelectStmt
		// ForStmt
		// RangeStmt

		switch {
		case funcLitOK:
			logger.S.Infof("funcLit %d --> %d", funcLitExpr.Body.Lbrace, funcLitExpr.Body.Rbrace)
			logger.S.Info(string(f.Content[funcLitExpr.Body.Lbrace-1 : funcLitExpr.Body.Rbrace-1]))
		case blockStmtOK:
			logger.S.Infof("block %d --> %d", blockStmt.Lbrace, blockStmt.Rbrace)
			logger.S.Info(string(f.Content[blockStmt.Lbrace-1 : blockStmt.Rbrace-1]))
		case assignOK:
			logger.S.Infof("assign %s", assignStmt.Tok.String())
		case identOK:
			logger.S.Info(identExpr.Name)
			logger.S.Infof("%#v", identExpr)
			if identExpr.Obj != nil {
				logger.S.Infof("%#v", identExpr.Obj)
				if identExpr.Obj.Decl != nil {
					logger.S.Infof("%#v", identExpr.Obj.Decl)
				}
			}

		case callOK:
			if !DetectInstanceCall(callExpr.Fun, platform.GinInstanceCall) {
				return false
			}

			index := f.NodeIndex(node)
			// idenNode := f.FindCallLIdent(index)
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
					switch sel.Name {
					case "GET", "PUT", "PATH", "POST", "DELETE":
						logger.S.Infof("x %#v", x)
						logger.S.Infof("sel %#v", sel)
						logger.S.Infof("args %#v", callExpr.Args)
						logger.S.Infof("args %#v", callExpr.Args[0]) // path
						if len(callExpr.Args) > 1 {
							logger.S.Infof("args %#v", callExpr.Args[1]) // func
						}

						p := ""
						if bl, blOK := callExpr.Args[0].(*ast.BasicLit); blOK {
							p = bl.Value
						}

						APIs = append(APIs, API{
							f,
							p,
							sel.Name,
						})

					}
				}

			}
		}
		return true
	})
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
		switch {
		case identOK:
			if ident.Obj != nil {
				if ident.Obj.Kind == ast.Var {
					if assignStmt, assignStmtOK := ident.Obj.Decl.(*ast.AssignStmt); assignStmtOK {
						if len(assignStmt.Rhs) != 0 {
							if callExpr, callExprOK := assignStmt.Rhs[0].(*ast.CallExpr); callExprOK {
								if selectorExpr, selectorExprOK := callExpr.Fun.(*ast.SelectorExpr); selectorExprOK {
									if equalSelectorExpr(selectorExpr, "gin", "Default") {
										f.GinIdents = append(f.GinIdents, ident)
									}
								}

							}
						}
					}

				}
			}

		case callExprOK:
			_ = callExpr
			// if selectorExpr, selectorExprOK := callExpr.Fun.(*ast.SelectorExpr); selectorExprOK {
			// 	switch selectorExpr.Sel.Name {
			// 	case `GET`, `POST`, `PUT`, `PATCH`, `DELETE`:
			// 		logger.S.Infof("callExpr %#v", callExpr)
			// 		logger.S.Infof("callExpr.Fun %#v", callExpr.Fun)
			// 		logger.S.Infof("callExpr.selectorExpr %#v", selectorExpr)
			// 		logger.S.Infof("callExpr.selectorExpr.X %#v", selectorExpr.X)
			// 		logger.S.Infof("callExpr.selectorExpr.Sel %#v", selectorExpr.Sel)

			// 		if iden, idenOK := selectorExpr.X.(*ast.Ident); idenOK {
			// 			logger.S.Infof("callExpr.selectorExpr.X.Obj %#v", iden.Obj)
			// 			logger.S.Infof("callExpr.selectorExpr.X.Obj.Name %#v", iden.Name)

			// 		}
			// 	}
			// }
		case typeSpecOK:
			typePos := -1
			if typeSpec.Type.Pos().IsValid() {
				typePos = int(typeSpec.Type.Pos())
			}
			Structs = append(Structs, Struct{f.FullPath, typeSpec.Name.Name, typePos, int(typeSpec.Name.NamePos), f.AstFile, &n})

		case structTypeOK:
			_ = structType
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
		logger.S.Infof("%#v", lNode)
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
