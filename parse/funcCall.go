package parse

import "go/ast"

// InstanceCall ...
type InstanceCall struct {
	Pkg      string
	FuncName string
}

var ginInstanceCall = []InstanceCall{
	{`gin`, `Default`},
	{`gin`, `New`},
}

func DetectInstanceCall(expr ast.Expr, calls []InstanceCall) bool {
	for _, v := range calls {
		if isFuncCallWithName(expr, v.Pkg, v.FuncName) {
			return true
		}
	}
	return false
}

func isFuncCallWithName(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}
