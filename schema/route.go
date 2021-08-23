package schema

import "go/ast"

type APIResopnse struct {
	Code string
}

type APIParam struct {
	StructName string
}

type GinAPI struct {
	Source      *SourceFile
	Group       string
	Method      string
	Path        string
	Node        *ast.Node
	APIParam    *APIParam
	APIResopnse []*APIResopnse
}

func NewGinRoute() *GinAPI {
	return &GinAPI{}
}
