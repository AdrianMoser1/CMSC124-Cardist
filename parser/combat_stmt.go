package parser

import "cardist/scanner"

type TurnStmt struct {
	Keyword scanner.Token
	Body    []Stmt
}

func (*TurnStmt) stmtNode() {}

func (p *Parser) turnStatement() Stmt {
	keyword := p.previous()
	if _, err := p.consume(scanner.TOKEN_LEFT_BRACE, "Expect '{' after 'turn'."); err != nil {
		return nil
	}
	return &TurnStmt{Keyword: keyword, Body: p.block()}
}
