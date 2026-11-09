package parser

import "cardist/scanner"

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (*WhileStmt) stmtNode() {}

// whileStatement parses after `while`: "(" expression ")" statement.
func (p *Parser) whileStatement() Stmt {
	if _, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'while'."); err != nil {
		return nil
	}
	condition := p.expression()
	if p.panicMode {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after condition."); err != nil {
		return nil
	}
	body := p.statement()
	if p.panicMode {
		return nil
	}
	return &WhileStmt{Condition: condition, Body: body}
}
