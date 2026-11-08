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

// BlockStmt is `{ declaration* }`. Its children may themselves be blocks,
// which is what makes scopes nest to any depth.
type BlockStmt struct {
	Statements []Stmt
}

func (*BlockStmt) stmtNode() {}
