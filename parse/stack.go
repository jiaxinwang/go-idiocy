package parse

import (
	"go/ast"
	"idiocy/logger"
)

func (f *SourceCode) BuildStacks() {
	f.fullStacks = []ast.Node{}
	f.walk(func(node ast.Node) bool {
		if node == nil {
			return false
		}
		f.fullStacks = append(f.fullStacks, node)
		return true
	})
}

func (f *SourceCode) NodeIndex(node ast.Node) int {
	for k, v := range f.fullStacks {
		switch {
		case v.Pos().IsValid() != node.Pos().IsValid():
			fallthrough
		case v.End().IsValid() != node.End().IsValid():
			continue
		}

		if node.Pos().IsValid() {
			if int(v.Pos()) != int(node.Pos()) {
				continue
			}
		}
		if node.End().IsValid() {
			if int(v.End()) != int(node.End()) {
				continue
			}
		}
		return k
	}
	return -1
}

func (f *SourceCode) StacksLength() int {
	return len(f.fullStacks)
}

// /ast.Ident
func (f *SourceCode) FindCallLIdent(callIndex int) ast.Node {
	lIndex := callIndex - 1
	llIndex := callIndex - 2
	if lIndex < 0 || llIndex < 0 {
		return nil
	}

	lNode := f.fullStacks[lIndex]
	llNode := f.fullStacks[llIndex]

	_, identOK := lNode.(*ast.Ident)                 //ObjKind.var
	assignStmt, assignOK := llNode.(*ast.AssignStmt) // Lhs,IsOperator
	if !assignOK || !identOK {
		return nil
	}

	if !assignStmt.Tok.IsOperator() {
		return nil
	}

	switch assignStmt.Tok.String() {
	case "=", ":=":
		logger.S.Infof("%#v", lNode)
		return lNode
	}

	return nil
}
