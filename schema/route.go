package schema

import "go/ast"

type GinRoute struct {
	Source *SourceFile
	Method string
	Path   string
	Node   *ast.Node
}

func NewGinRoute() *GinRoute {
	return &GinRoute{}
}
