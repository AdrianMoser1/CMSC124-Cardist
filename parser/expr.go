package parser

import "cardist/scanner"

type Expr interface {
	Accept(visitor Visitor) any
}

// changed method signatures to any so it can also return numbers, bool ,str, etc.
type Visitor interface {
	VisitLiteral(expr *LiteralExpr) any
	VisitUnary(expr *UnaryExpr) any
	VisitBinary(expr *BinaryExpr) any
	VisitGrouping(expr *GroupingExpr) any
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v Visitor) any {
	return v.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v Visitor) any {
	return v.VisitUnary(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v Visitor) any {
	return v.VisitBinary(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v Visitor) any {
	return v.VisitGrouping(e)
}
