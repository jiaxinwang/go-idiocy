package schema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"idiocy/logger"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type SourceFile struct {
	Name    string
	Path    string
	PkgPath string
	Imports []*ast.ImportSpec
}

func (f *SourceFile) ParseFile(filename string) error {
	fset := token.NewFileSet()
	path, _ := filepath.Abs(filename)
	astFile, err := parser.ParseFile(fset, path, nil, parser.AllErrors)
	if err != nil {
		logrus.WithError(err).Error()
		return err
	}

	f.Name = astFile.Name.Name
	f.Imports = astFile.Imports
	logger.S.Info(astFile.Scope.Objects)
	logger.S.Info(astFile.Comments)
	logger.S.Info(astFile.Name.Obj)

	for _, v := range astFile.Imports {
		logger.S.Infow("import", "path", v.Path)
	}

	for _, node := range astFile.Decls {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					switch typeSpec.Type.(type) {
					case *ast.StructType:
						structType := typeSpec.Type.(*ast.StructType)
						for _, field := range structType.Fields.List {
							ident, ok := field.Type.(*ast.Ident)
							if ok {
								fieldType := ident.Name
								for _, name := range field.Names {
									logrus.WithField("type", "*ast.Ident").WithField(name.Name, fieldType).Info()
								}
								continue
							}
							selectorExpr, ok := field.Type.(*ast.SelectorExpr)
							if ok {
								fieldType := selectorExpr.Sel.Name
								for _, name := range field.Names {
									logrus.WithField("type", "*ast.SelectorExpr").WithField(name.Name, fieldType).Info()
								}
								continue

							}

							starExpr, ok := field.Type.(*ast.StarExpr)
							if ok {
								selectorExpr, ok := starExpr.X.(*ast.SelectorExpr)
								if !ok {
									continue
								}
								fieldType := selectorExpr.Sel.Name
								for _, name := range field.Names {
									logrus.WithField("type", "*ast.StarExpr").WithField(name.Name, fieldType).Info()
								}
								continue
							}
						}
					}
				}
			}
		}
	}
	return nil
}

type v struct {
	info *types.Info
}

func (v v) Visit(node ast.Node) (w ast.Visitor) {
	switch node := node.(type) {
	case *ast.CallExpr:
		switch node := node.Fun.(type) {
		case *ast.SelectorExpr: // foo.ReadFile
			pkgID := node.X.(*ast.Ident)
			fmt.Println(v.info.Uses[pkgID].(*types.PkgName).Imported().Path())
		case *ast.Ident: // ReadFile
			pkgID := node
			fmt.Println(v.info.Uses[pkgID].Pkg().Path())
		}
	}

	return v
}
