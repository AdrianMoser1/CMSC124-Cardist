package parser

import "cardist/scanner"

// Program is the root node of the AST. Its children are the file's
// expressions in source order (an n-ary node: children live in a slice, the
// same way any n-ary tree node stores them). Because Program is itself an
// Expr, the whole file is ONE tree that a visitor can walk from the root.
type Program struct {
	Expressions []Expr
}

func (p *Program) Accept(v Visitor) any {
	return v.VisitProgram(p)
}

type Expr interface {
	Accept(visitor Visitor) any
}

// changed method signatures to any so it can also return numbers, bool ,str, etc.
type Visitor interface {
	VisitProgram(program *Program) any
	VisitLiteral(expr *LiteralExpr) any
	VisitUnary(expr *UnaryExpr) any
	VisitBinary(expr *BinaryExpr) any
	VisitGrouping(expr *GroupingExpr) any
	VisitVariable(expr *VariableExpr) any
	VisitAssign(expr *AssignExpr) any
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

type VariableExpr struct {
	Name scanner.Token
}

func (e *VariableExpr) Accept(v Visitor) any {
	return v.VisitVariable(e)
}

// AssignExpr writes to an existing binding. Assignment is an expression,
// so `var b = (a = 1);` works.
type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

func (e *AssignExpr) Accept(v Visitor) any {
	return v.VisitAssign(e)
}
