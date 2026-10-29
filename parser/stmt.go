package parser


type Stmt interface { 
	stmtNode()
}

//root of a statement-based program
type Script struct {
	Statements []Stmt
}

type PrintStmt struct {
	Expression Expr
}

func (*PrintStmt) stmtNode() {}

type ExpressionStmt struct{
	Expression Expr
}

func (*ExpressionStmt) stmtNode {}

