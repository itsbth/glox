package main

import (
	"fmt"

	"github.com/itsbth/glox/ast"
)

// PrettyPrinter dummy
type PrettyPrinter struct{}

func (p PrettyPrinter) VisitBinary(node *ast.Binary) interface{} {
	return fmt.Sprintf("(%s %s %s)", node.Operator(), node.Left().Accept(p), node.Right().Accept(p))
}

func (p PrettyPrinter) VisitGrouping(node *ast.Grouping) interface{} {
	return fmt.Sprintf("(%s)", node.Expr().Accept(p))
}

func (p PrettyPrinter) VisitLiteral(node *ast.Literal) interface{} {
	return node.Value()
}

func (p PrettyPrinter) VisitUnary(node *ast.Unary) interface{} {
	return fmt.Sprintf("(%s %s)", node.Operator(), node.Right().Accept(p))
}

func (p PrettyPrinter) VisitIdentifier(node *ast.Identifier) interface{} {
	return node.Name()
}
