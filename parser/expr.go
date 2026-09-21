package parser

import "cardist/scanner"

type Expr interface {
	Accept(visitor Visitor) string
}

type Visitor interface {
	VisitLiteral(expr *LiteralExpr) string
	VisitUnary(expr *UnaryExpr) string
	VisitBinary(expr *BinaryExpr) string
	VisitGrouping(expr *GroupingExpr) string
}

type LiteralExpr struct {
	Value any
}

func (e *LiteralExpr) Accept(v Visitor) string {
	return v.VisitLiteral(e)
}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (e *UnaryExpr) Accept(v Visitor) string {
	return v.VisitUnary(e)
}

type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

func (e *BinaryExpr) Accept(v Visitor) string {
	return v.VisitBinary(e)
}

type GroupingExpr struct {
	Expression Expr
}

func (e *GroupingExpr) Accept(v Visitor) string {
	return v.VisitGrouping(e)
}
