package schema

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/jiaxinwang/common/fs"
	"github.com/jiaxinwang/err2"
)

type API struct {
	SourceFile *SourceFile
	Path       string
	Method     string
}

type SourceFile struct {
	FullPath       string
	Path           string
	PkgPath        string
	Imports        []*ast.ImportSpec
	FileSet        *token.FileSet
	AstFile        *ast.File
	Content        []byte
	main           bool
	Decls          map[string]string
	fullStacks     []ast.Node
	GinIdents      []*ast.Ident
	GinIdentifiers []*GinIdentifier
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
