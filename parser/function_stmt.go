package parser

import "cardist/scanner"

// FunctionStmt is `func name(params) { body }`.
type FunctionStmt struct {
	Name   scanner.Token
	Params []scanner.Token
	Body   []Stmt
}

func (*FunctionStmt) stmtNode() {}

// funDeclaration parses after `func`: IDENTIFIER "(" parameters? ")" block.
func (p *Parser) funDeclaration() Stmt {
	name, err := p.consume(scanner.TOKEN_IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after function name."); err != nil {
		return nil
	}

	var params []scanner.Token
	if !p.check(scanner.TOKEN_RIGHT_PAREN) {
		for {
			param, err := p.consume(scanner.TOKEN_IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil
			}
			params = append(params, param)
			if !p.match(scanner.TOKEN_COMMA) {
				break
			}
		}
	}

	if _, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after parameters."); err != nil {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_LEFT_BRACE, "Expect '{' before function body."); err != nil {
		return nil
	}
	body := p.block()
	if p.panicMode {
		return nil
	}
	return &FunctionStmt{Name: name, Params: params, Body: body}
}
