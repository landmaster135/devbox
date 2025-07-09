// internal/analyzer/complexity.go
package analyzer

import (
	"go/ast"
	"go/token"
)

// 複雑度計算用ビジター
type complexityVisitor struct {
	complexity int
}

// Visit implements ast.Visitor for calculating cyclomatic complexity
func (v *complexityVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause, *ast.BinaryExpr:
		if bin, ok := n.(*ast.BinaryExpr); ok {
			if bin.Op == token.LAND || bin.Op == token.LOR {
				v.complexity++
			}
		} else {
			v.complexity++
		}
	}
	return v
}

// 関数分析用ビジター
type funcVisitor struct {
	count     int
	sizes     []int
	maxSize   int
	totalSize int
}

// Visit implements ast.Visitor for analyzing functions
func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	if fn, ok := node.(*ast.FuncDecl); ok {
		v.count++

		fset := token.NewFileSet()
		size := fset.Position(fn.End()).Line - fset.Position(fn.Pos()).Line

		v.sizes = append(v.sizes, size)
		v.totalSize += size

		if size > v.maxSize {
			v.maxSize = size
		}
	}
	return v
}
