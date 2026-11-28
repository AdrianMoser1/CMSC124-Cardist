package parser

import "cardist/scanner"

type ControlStmt struct {
	Keyword scanner.Token
}

func (*ControlStmt) stmtNode() {}

func (p *Parser) controlStatement() Stmt {
	keyword := p.previous()
	if keyword.Type == scanner.TOKEN_BREAK && p.loopDepth == 0 && p.switchDepth == 0 {
		p.reportError(keyword, "Can't use 'break' outside a loop or switch.")
		return nil
	}
	if keyword.Type == scanner.TOKEN_CONTINUE && p.loopDepth == 0 {
		p.reportError(keyword, "Can't use 'continue' outside a loop.")
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_SEMI_COLON, "Expect ';' after '"+keyword.Lexeme+"'."); err != nil {
		return nil
	}
	return &ControlStmt{Keyword: keyword}
}
