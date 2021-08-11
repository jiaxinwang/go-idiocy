package schema

import (
	"go/ast"
	"go/parser"
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
		// only interested in generic declarations
		if genDecl, ok := d.(*ast.GenDecl); ok {

			// handle const's and vars
			if genDecl.Tok == token.CONST || genDecl.Tok == token.VAR {

				// there may be multiple
				// i.e. const ( ... )
				for _, cDecl := range genDecl.Specs {

					// havn't find another kind of spec then value but better check
					if vSpec, ok := cDecl.(*ast.ValueSpec); ok {
						log.Printf("const ValueSpec: %+v\n", vSpec)

						// iterate over Name/Value pair
						for i := 0; i < len(vSpec.Names); i++ {
							// TODO: only basic literals work currently
							switch v := vSpec.Values[i].(type) {
							case *ast.BasicLit:
								f.Decls[vSpec.Names[i].Name] = v.Value
							default:
								log.Printf("Name: %s - Unsupported ValueSpec: %+v\n", vSpec.Names[i].Name, v)
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

		// switch x := callExpr.Args[0].(type) {
		// case *ast.BasicLit:
		// 	log.Println("Literal Argument:", x.Value)
		// 	f.boxes = append(f.boxes, Box{x.Value})

		// case *ast.Ident:
		// 	log.Printf("Argument Identifier: %+v", x)
		// 	val, ok := f.decls[x.Name]
		// 	if !ok {
		// 		//TODO: Add ERRORs list to file type and return after iteration!
		// 		log.Printf("Could not find identifier[%s] in decls map\n", x.Name)
		// 		return true
		// 	}
		// 	f.boxes = append(f.boxes, Box{val})

		// default:
		// 	fmt.Println("Unsupported argument to rice.(must)FindBox():", x)
		// }

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
type Box struct {
	name string
}

func FindRiceBoxes(filename string, src []byte) error {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return err
	}

	f := &SourceFile{FileSet: fset, AstFile: astFile, Content: src}
	f.Decls = make(map[string]string)
	f.find()
	return nil

}
func (f *SourceFile) EnumerateStruct() {
	ast.Inspect(f.AstFile, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		nodeIndex := f.NodeIndex(n)
		_ = nodeIndex
		structType, structTypeOK := n.(*ast.StructType)
		typeSpec, typeSpecOK := n.(*ast.TypeSpec)
		switch {
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
