package helper

import (
	"go/ast"
)

//
func ExtractExprStmt(node ast.Node) (x *ast.Expr, ok bool) {
	exprStmt, exprStmtOK := node.(*ast.ExprStmt)
	if exprStmtOK {
		return &exprStmt.X, true
	}
	return nil, false
}

func ExtractCallExpr(node ast.Node) (fun *ast.Expr, args []ast.Expr, ok bool) {
	// callExpr.Args
	// callExpr.Ellipsis
	// callExpr.Fun
	// callExpr.Lparen
	// callExpr.Rparen

	callExpr, callExprOK := node.(*ast.CallExpr)
	if callExprOK {
		return &callExpr.Fun, callExpr.Args, true
	}
	return nil, []ast.Expr{}, false

}

func ExtractSelectorExpr(node ast.Node) (x *ast.Expr, sel *ast.Ident, ok bool) {
	selectorExpr, selectorExprOK := node.(*ast.SelectorExpr)
	if selectorExprOK {
		return &selectorExpr.X, selectorExpr.Sel, true
	}
	return nil, nil, false
}

func ExtractBasicLit(node ast.Node) (kind int, value string, ok bool) {
	basicLit, basicLitOK := node.(*ast.BasicLit)
	if basicLitOK {
		return int(basicLit.Kind), basicLit.Value, true
	}
	return -1, "", false
}

func ExtractIdent(node ast.Node) (name string, object *ast.Object, ok bool) {
	ident, identOK := node.(*ast.Ident)
	if identOK {
		return ident.Name, ident.Obj, true
	}
	return "", nil, false
}
