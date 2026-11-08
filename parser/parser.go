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

	// panicMode is set when a syntax error is reported and cleared once
	// synchronize() has found a point to resume parsing. Unlike hadError (which
	// stays true for the whole file so main can exit 65), it says "the
	// statement being parsed right now has already failed".
	panicMode bool
	// blockDepth is how many '{' are open. synchronize() uses it to leave a
	// '}' for the enclosing block instead of swallowing it.
	blockDepth int
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

func (p *Parser) Parse() *Program {
	program := &Program{}
	for !p.isAtEnd() {
		earlier := p.hadError
		p.hadError = false // track whether *this* expression failed, not the whole file

		expr := p.expression()

		// Optional semicolon check so standalone test expressions don't throw syntax errors
		p.match(scanner.TOKEN_SEMI_COLON)

		if p.hadError {
			p.synchronize()
		} else if expr != nil {
			program.Expressions = append(program.Expressions, expr)
		}
		p.hadError = p.hadError || earlier
	}
	return program
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

	// Added: Match identifiers (variable names) and store their lexeme string
	if p.match(scanner.TOKEN_IDENTIFIER) {
		return &VariableExpr{Name: p.previous()}
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
	p.panicMode = true
	if token.Type == scanner.TOKEN_EOF {
		fmt.Fprintf(os.Stderr, "[line %d] Error at end: %s\n", token.Line, message)
	} else {
		fmt.Fprintf(os.Stderr, "[line %d] Error at '%s': %s\n", token.Line, token.Lexeme, message)
	}
}

// synchronize discards tokens after a syntax error until a new statement can
// plausibly begin: just after a ';', or just before a statement keyword or '{'.
//
// It always consumes at least one token, so a caller that loops "parse a
// statement, synchronize on error" makes progress and can never spin. The one
// exception is a '}' that closes an open block: that is left for block(),
// otherwise recovery inside a block would eat its own closing brace and report
// a bogus "Expect '}'" afterwards.
func (p *Parser) synchronize() {
	p.panicMode = false

	if p.isAtEnd() {
		return
	}
	if p.blockDepth > 0 && p.check(scanner.TOKEN_RIGHT_BRACE) {
		return
	}
	p.advance() // guaranteed progress

	for !p.isAtEnd() {
		if p.previous().Type == scanner.TOKEN_SEMI_COLON {
			return
		}

		switch p.peek().Type {
		case scanner.TOKEN_VAR, scanner.TOKEN_IF, scanner.TOKEN_ELSEIF, scanner.TOKEN_ELSE,
			scanner.TOKEN_SWITCH, scanner.TOKEN_WHILE, scanner.TOKEN_PRINT, scanner.TOKEN_RETURN,
			scanner.TOKEN_FOR, scanner.TOKEN_FUNC, scanner.TOKEN_CARD, scanner.TOKEN_ENEMY,
			scanner.TOKEN_ARTIFACT, scanner.TOKEN_ELIXIR, scanner.TOKEN_TURN,
			scanner.TOKEN_LEFT_BRACE:
			return
		case scanner.TOKEN_RIGHT_BRACE:
			if p.blockDepth > 0 {
				return
			}
		}

		p.advance()
	}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment parses the left side as an ordinary expression FIRST, then
// checks whether it was a legal target. The parser can't know `a + b` is
// wrong until it sees the `=` after it, so it validates afterwards.
//
// "Has this statement already failed?" is panicMode, not hadError: hadError
// stays true for the rest of the file once ANY earlier statement failed, which
// would make every later assignment look like a failure and cascade into
// bogus errors.
func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.panicMode || !p.match(scanner.TOKEN_EQUAL) {
		return expr
	}

	equals := p.previous()
	value := p.assignment() // recursing on the right makes `a = b = 1` right-associative
	if p.panicMode {
		return nil
	}

	if target, ok := expr.(*VariableExpr); ok {
		return &AssignExpr{Name: target.Name, Value: value}
	}

	p.reportError(equals, "Invalid assignment target.")
	return nil

}
