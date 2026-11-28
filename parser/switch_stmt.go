package parser

import "cardist/scanner"

type SwitchStmt struct {
	Subject Expr
	Cases   []SwitchCase
	Default *BlockStmt
}

type SwitchCase struct {
	Value Expr
	Body  *BlockStmt
}

func (*SwitchStmt) stmtNode() {}

func (p *Parser) switchStatement() Stmt {
	if _, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'switch'."); err != nil {
		return nil
	}
	subject := p.expression()
	if p.panicMode {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after switch expression."); err != nil {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_LEFT_BRACE, "Expect '{' before switch cases."); err != nil {
		return nil
	}

	p.switchDepth++
	p.blockDepth++
	defer func() {
		p.switchDepth--
		p.blockDepth--
	}()

	stmt := &SwitchStmt{Subject: subject}
	for !p.check(scanner.TOKEN_RIGHT_BRACE) && !p.isAtEnd() {
		switch {
		case p.match(scanner.TOKEN_CASE):
			value := p.expression()
			if p.panicMode {
				return nil
			}
			body := p.switchBody()
			if body == nil {
				return nil
			}
			stmt.Cases = append(stmt.Cases, SwitchCase{Value: value, Body: body})
		case p.match(scanner.TOKEN_DEFAULT):
			if stmt.Default != nil {
				p.reportError(p.previous(), "A switch can have only one default.")
				return nil
			}
			stmt.Default = p.switchBody()
			if stmt.Default == nil {
				return nil
			}
		default:
			p.reportError(p.peek(), "Expect 'case' or 'default' in switch.")
			return nil
		}
	}
	if _, err := p.consume(scanner.TOKEN_RIGHT_BRACE, "Expect '}' after switch cases."); err != nil {
		return nil
	}
	return stmt
}

func (p *Parser) switchBody() *BlockStmt {
	if _, err := p.consume(scanner.TOKEN_COLON, "Expect ':' after switch label."); err != nil {
		return nil
	}

	body := &BlockStmt{}
	for !p.isAtEnd() && !p.check(scanner.TOKEN_RIGHT_BRACE) &&
		!p.check(scanner.TOKEN_CASE) && !p.check(scanner.TOKEN_DEFAULT) {
		if stmt := p.declaration(); stmt != nil {
			body.Statements = append(body.Statements, stmt)
		}
	}
	return body
}
