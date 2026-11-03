package parser

import (
	"cardist/scanner"
)

// VarStmt declares a name. Initializer is nil for `var x;`.
type VarStmt struct {
	Name        scanner.Token
	Initializer Expr
}

func (*VarStmt) stmtNode() {}

type Stmt interface {
	stmtNode()
}

// root of a statement-based program
type Script struct {
	Statements []Stmt
}

type PrintStmt struct {
	Expression Expr
}

func (*PrintStmt) stmtNode() {}

type ExpressionStmt struct {
	Expression Expr
}

func (*ExpressionStmt) stmtNode() {}
