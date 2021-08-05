package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"idiocy/logger"
	"log"
	"os"
	"path/filepath"

	"github.com/jiaxinwang/common/fs"
	"github.com/jiaxinwang/err2"
)

type sourceCode struct {
	fset       *token.FileSet
	astFile    *ast.File
	src        []byte
	filename   string
	main       bool
	decls      map[string]string
	boxes      []Box
	fullStacks []ast.Node
}

func (f *sourceCode) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *sourceCode) find() {
	f.BuildStacks()
	f.findDecals()
	f.findGinInstance()
}

// Gin ...
func Gin(dir string) error {
	var err error
	defer err2.Return(&err)
	err = filepath.Walk(dir,
		func(curDir string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fset := token.NewFileSet()
				path, _ := filepath.Abs(curDir)
				f, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
				if err != nil {
					return err
				}

				b := err2.Bytes.Try(fs.ReadBytes(curDir))

				code := &sourceCode{
					fset:     fset,
					astFile:  f,
					src:      b,
					filename: curDir,
					decls:    make(map[string]string),
				}
				code.find()

			}
			return err
		})
	if err != nil {
		return err
	}
	return nil
}

func (f *sourceCode) findDecals() {
	for _, d := range f.astFile.Decls {
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
								f.decls[vSpec.Names[i].Name] = v.Value
							default:
								log.Printf("Name: %s - Unsupported ValueSpec: %+v\n", vSpec.Names[i].Name, v)
							}
						}
					}
				}
			}

		}
	}

	log.Println("Decls:", f.decls)
}

func (f *sourceCode) findGinInstance() {

	// var lastAssignStmt *ast.AssignStmt
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}

		logger.S.Infof("----- %d --> %d -----", int(node.Pos()), int(node.End()))
		logger.S.Info(string(f.src[node.Pos()-1 : node.End()-1]))

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
			logger.S.Info(string(f.src[funcLitExpr.Body.Lbrace-1 : funcLitExpr.Body.Rbrace-1]))
		case blockStmtOK:
			logger.S.Infof("block %d --> %d", blockStmt.Lbrace, blockStmt.Rbrace)
			logger.S.Info(string(f.src[blockStmt.Lbrace-1 : blockStmt.Rbrace-1]))
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
			if !DetectInstanceCall(callExpr.Fun, ginInstanceCall) {
				return false
			}

			index := f.NodeIndex(node)
			idenNode := f.FindCallLIdent(index)

		default:
			logger.S.Infof(logger.ColorLightGreen("unhandle %d --> %d"), node.Pos(), node.End())
			logger.S.Info(logger.ColorLightGreen(string(f.src[node.Pos()-1 : node.End()-1])))
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

	f := &sourceCode{fset: fset, astFile: astFile, src: src, filename: filename}
	f.decls = make(map[string]string)
	f.find()
	return nil

}
