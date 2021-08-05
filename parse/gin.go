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
	fset     *token.FileSet
	astFile  *ast.File
	src      []byte
	filename string
	main     bool
	decls    map[string]string
	boxes    []Box
}

func (f *sourceCode) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *sourceCode) find() {
	f.findDecals()
	f.findGinDefault()
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

func (f *sourceCode) findGinDefault() {
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}

		// logger.S.Infof("-----")
		// logger.S.Infof("%#v", node)
		// logger.S.Infof("%#v", node.Pos())
		if node.Pos().IsValid() {
			logger.S.Infof("node.Pos() %d", int(node.Pos()))
		}
		if node.End().IsValid() {
			logger.S.Infof("node.End() %d", int(node.End()))
		}
		b := f.src
		// f.src[9:13]
		logger.S.Info(len(b))
		// logger.S.Info(string(b[int(node.Pos()) : int(node.End())-1]))

		// logger.S.Infof("----------")
		identExpr, identOK := node.(*ast.Ident)
		callExpr, callOK := node.(*ast.CallExpr)
		assignExpr, assignOK := node.(*ast.AssignStmt)
		switch {
		case assignOK:
			logger.S.Info(assignExpr.Tok.String())
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
			isMustFindBox := isFuncCallWithName(callExpr.Fun, "gin", "Default")
			if !(isMustFindBox) || len(callExpr.Args) != 0 {
				return false
			}
			// logger.S.Infof("%#v", callExpr)
			// logger.S.Infof("%+v", callExpr.Args)
			// logger.S.Infof("%+v", callExpr.Ellipsis)
			// logger.S.Infof("%+v", callExpr.Fun)
			// logger.S.Infof("%+v", callExpr.Lparen)
			// logger.S.Infof("%+v", callExpr.Rparen)

			// logger.S.Infow(fmt.Sprintf("gin.Default Call! @%+v", callExpr),
			// 	"pos", node.Pos(),
			// 	"end", node.End(),
			// )
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
func isFuncCallWithName(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

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
