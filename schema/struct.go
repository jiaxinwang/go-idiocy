package schema

import "go/ast"

var Structs = []Struct{}

type Struct struct {
	FullPath string
	Name     string
	TypePos  int
	NamePos  int
	AstFile  *ast.File
	Node     *ast.Node
}
