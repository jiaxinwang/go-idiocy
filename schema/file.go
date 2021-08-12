package schema

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/jiaxinwang/common/fs"
	"github.com/jiaxinwang/err2"
)

type API struct {
	SourceFile *SourceFile
	Path       string
	Method     string
}

type SourceFile struct {
	FullPath   string
	Path       string
	PkgPath    string
	Imports    []*ast.ImportSpec
	FileSet    *token.FileSet
	AstFile    *ast.File
	Content    []byte
	main       bool
	Decls      map[string]string
	fullStacks []ast.Node
	GinIdents  []*ast.Ident
}

func (f *SourceFile) ParseFile() error {

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, f.FullPath, nil, parser.AllErrors)
	if err != nil {
		return err
	}

	b := err2.Bytes.Try(fs.ReadBytes(f.FullPath))

	f.FileSet = fset
	f.AstFile = astFile
	f.Content = b
	f.Decls = make(map[string]string)

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
