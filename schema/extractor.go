package schema

import (
	"go/ast"
	"strings"
)

func extractSelectorExpr(node *ast.SelectorExpr) (X *ast.Ident, Sel *ast.Ident, ok bool) {
	if X, ok = node.X.(*ast.Ident); !ok {
		return nil, nil, false
	}
	Sel = node.Sel
	return
}

func equalSelectorExpr(selectorExpr *ast.SelectorExpr, xName, selName string) bool {
	x, sel, ok := extractSelectorExpr(selectorExpr)
	if !ok {
		return false
	}
	return equalIdent(x, xName) && equalIdent(sel, selName)
}

func equalIdent(ident *ast.Ident, name string) bool {
	return strings.EqualFold(ident.Name, name)
}
