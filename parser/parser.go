package parser

import (
	"cardist/scanner"
	"fmt"
	"os"
)

type Parser struct {
	tokens   []scanner.Token
	current  int
	hadError bool
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		current:  0,
		hadError: false,
	}
}

func (p *Parser) ErrorFound() bool {
	return p.hadError
}

func (p *Parser) Parse() []Expr {
	var expressions []Expr
	for !p.isAtEnd() {
		earlier := p.hadError
		p.hadError = false // track whether *this* expression failed, not the whole file

		expr := p.expression()
		if !p.hadError {
			p.consume(scanner.TOKEN_SEMI_COLON, "Expect ';' after expression.")
		}
		if p.hadError {
			p.synchronize()
		} else if expr != nil {
			expressions = append(expressions, expr)
		}
		p.hadError = p.hadError || earlier
	}
	return expressions
}

func (p *Parser) expression() Expr {
	return p.or()
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(scanner.TOKEN_OR) {
		operator := p.previous()
		right := p.and()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(scanner.TOKEN_AND) {
		operator := p.previous()
		right := p.equality()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(scanner.TOKEN_BANG_EQUAL, scanner.TOKEN_EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(scanner.TOKEN_GREATER, scanner.TOKEN_GREATER_EQUAL, scanner.TOKEN_LESSER, scanner.TOKEN_LESSER_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(scanner.TOKEN_MINUS, scanner.TOKEN_PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(scanner.TOKEN_SLASH, scanner.TOKEN_STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(scanner.TOKEN_BANG, scanner.TOKEN_MINUS) {
		operator := p.previous()
		right := p.unary()
		return &UnaryExpr{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(scanner.TOKEN_FALSE) {
		return &LiteralExpr{Value: false}
	}
	if p.match(scanner.TOKEN_TRUE) {
		return &LiteralExpr{Value: true}
	}
	if p.match(scanner.TOKEN_NIL) {
		return &LiteralExpr{Value: nil}
	}

	if p.match(scanner.TOKEN_NUMBER, scanner.TOKEN_STRING) {
		return &LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(scanner.TOKEN_IDENTIFIER) {
		return &LiteralExpr{Value: p.previous().Lexeme}
	}

	if p.match(scanner.TOKEN_LEFT_PAREN) {
		expr := p.expression()
		_, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil
		}
		return &GroupingExpr{Expression: expr}
	}

	p.reportError(p.peek(), "Expect expression.")
	return nil
}

// Token navigation helpers
func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.TOKEN_EOF
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	p.reportError(p.peek(), message)
	return scanner.Token{}, fmt.Errorf("parse error")
}

func (p *Parser) reportError(token scanner.Token, message string) {
	p.hadError = true
	if token.Type == scanner.TOKEN_EOF {
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", token.Line, message)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, message)
	}
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.TOKEN_SEMI_COLON {
			return
		}

		switch p.peek().Type {
		case scanner.TOKEN_VAR, scanner.TOKEN_IF, scanner.TOKEN_WHILE,
			scanner.TOKEN_PRINT, scanner.TOKEN_RETURN, scanner.TOKEN_FOR,
			scanner.TOKEN_FUNC, scanner.TOKEN_CARD:
			return
		}

		p.advance()
	}
}
