package helper

import (
	"fmt"
	"go/ast"

	"github.com/jiaxinwang/go-idiocy/logger"
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

func ExtractUnary(node ast.Node) (unaryExpr *ast.UnaryExpr, ok bool) {
	unary, unaryOK := node.(*ast.UnaryExpr)
	if unaryOK {
		return unary, true
	}
	return nil, false
}
func ExtractCompositeLit(node ast.Node) (cl *ast.CompositeLit, ok bool) {
	cl, ok = node.(*ast.CompositeLit)
	return
}

func ExtractIdent(node ast.Node) (name string, object *ast.Object, ok bool) {
	ident, identOK := node.(*ast.Ident)
	if identOK {
		return ident.Name, ident.Obj, true
	}
	return "", nil, false
}

func ExtractValueSpec(node ast.Node) (
	doc *ast.CommentGroup,
	names []*ast.Ident,
	typeNode *ast.Expr,
	values []ast.Expr,
	comment *ast.CommentGroup,
	ok bool) {
	valueSpec, valueSpecOK := node.(*ast.ValueSpec)
	if valueSpecOK {
		return valueSpec.Doc, valueSpec.Names, &valueSpec.Type, valueSpec.Values, valueSpec.Comment, true
	}
	return nil, []*ast.Ident{}, nil, []ast.Expr{}, nil, false
}

func ExplainObjectType(node ast.SelectorExpr) {
	logger.S.Infof("type %#v", node)
}

func ExtractObjectTypeName(object *ast.Object) (fullname string) {
	pkg, name := "", ""
	if ts, tsOK := object.Decl.(*ast.ValueSpec); tsOK {
		if _, _, typeNode, _, _, ok := ExtractValueSpec(ts); ok {
			typeNodeSelectorExpr, _ := ts.Type.(*ast.SelectorExpr)
			pkg, _, _ = ExtractIdent(typeNodeSelectorExpr.X)
			_, sel, _ := ExtractSelectorExpr(*typeNode)
			name = sel.Name
		}
	}
	return fmt.Sprintf("%s.%s", pkg, name)
}

func ExtractFuncLit(node ast.Node) (fl *ast.FuncLit, ok bool) {
	if fl, flOK := node.(*ast.FuncLit); flOK {
		_ = fl
		return fl, flOK
	}
	return nil, false
}

func ExtractBlockStml(node ast.Node) (bl *ast.BlockStmt, ok bool) {
	if blk, blkOK := node.(*ast.BlockStmt); blkOK {
		_ = blk
		return blk, blkOK
	}
	return nil, false
}
