package parser

import "cardist/scanner"

// ReturnStmt is `return expression? ;`. Value is nil for a bare `return;`.
type ReturnStmt struct {
	Keyword scanner.Token
	Value   Expr
}

func (*ReturnStmt) stmtNode() {}

// returnStatement parses after `return`.
func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(scanner.TOKEN_SEMI_COLON) {
		value = p.expression()
		if p.panicMode {
			return nil
		}
	}
	if _, err := p.consume(scanner.TOKEN_SEMI_COLON, "Expect ';' after return value."); err != nil {
		return nil
	}
	return &ReturnStmt{Keyword: keyword, Value: value}
}
