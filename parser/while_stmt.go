package parser

import "cardist/scanner"

type WhileStmt struct {
	Condition Expr
	Body      Stmt
	Increment Expr // Set by desugared for loops so continue still runs it.
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
	body := p.loopBody()
	if p.panicMode {
		return nil
	}
	return &WhileStmt{Condition: condition, Body: body}
}

func (p *Parser) loopBody() Stmt {
	p.loopDepth++
	defer func() { p.loopDepth-- }()
	return p.statement()
}
