package schema

import "go/ast"

type GinIdentifier struct {
	Source         *SourceFile
	InstancingCall *ast.CallExpr
	Node           *ast.Ident
	Calls          []*ast.CallExpr
}

func NewGinIdentifier() *GinIdentifier {
	return &GinIdentifier{
		Calls: []*ast.CallExpr{},
	}
}

func (g *GinIdentifier) Equal(another *GinIdentifier) bool {
	if !g.Source.Equal(another.Source) {
		return false
	}
	if g.InstancingCall.Pos() != another.InstancingCall.Pos() {
		return false
	}
	// if g.Node.Pos() != another.Node.Pos() {
	// 	return false
	// }
	// if g.Node.End() != another.Node.End() {
	// 	return false
	// }
	return true
}

func (g *GinIdentifier) AddCall(call *ast.CallExpr) bool {
	g.Calls = append(g.Calls, call)
	return true
}
