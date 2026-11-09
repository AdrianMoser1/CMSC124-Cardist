package parser

import "cardist/scanner"

// ParseScript parses statements for plain file execution.
// Parse() remains the expression-mode entry point for earlier labs.
//
// A syntax error does not stop the parse: declaration() reports it,
// synchronizes to the next statement boundary, and parsing continues, so every
// syntax error in the file is reported in one run. hadError stays true for the
// rest of the file, and the caller checks ErrorFound() and exits 65 before
// executing anything.
func (p *Parser) ParseScript() *Script {
	script := &Script{}

	for !p.isAtEnd() {
		if stmt := p.declaration(); stmt != nil {
			script.Statements = append(script.Statements, stmt)
		}
	}

	return script
}

// declaration is the recovery point. If the statement it just parsed reported
// an error (panicMode is still set), it resynchronizes and returns nil.
//
// Inside the statement parsers below, "did THIS statement fail?" is panicMode,
// not hadError. hadError stays true for the whole file once anything has
// failed, so testing it would discard every statement after the first error.
func (p *Parser) declaration() Stmt {
	start := p.current

	var stmt Stmt
	if p.match(scanner.TOKEN_VAR) {
		stmt = p.varDeclaration()
	} else {
		stmt = p.statement()
	}

	if p.panicMode {
		p.synchronize()
		// synchronize() always consumes a token, but make progress certain
		// here too: ParseScript and block() loop on declaration(), so a
		// declaration that consumed nothing would hang until the harness's
		// 15-second timeout.
		if p.current == start && !p.isAtEnd() {
			p.advance()
		}
		return nil
	}
	return stmt
}

func (p *Parser) varDeclaration() Stmt {
	name, err := p.consume(scanner.TOKEN_IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil
	}

	var initializer Expr // stays nil for `var x;`
	if p.match(scanner.TOKEN_EQUAL) {
		initializer = p.expression()
		if p.panicMode {
			return nil
		}
	}

	if _, err := p.consume(scanner.TOKEN_SEMI_COLON, "Expect ';' after variable declaration."); err != nil {
		return nil
	}
	return &VarStmt{Name: name, Initializer: initializer}
}

func (p *Parser) statement() Stmt {
	if p.match(scanner.TOKEN_IF) {
		return p.ifStatement()
	}
	if p.match(scanner.TOKEN_WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.TOKEN_PRINT) {
		return p.printStatement()
	}
	if p.match(scanner.TOKEN_LEFT_BRACE) {
		return &BlockStmt{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	if p.panicMode {
		return nil
	}

	_, err := p.consume(
		scanner.TOKEN_SEMI_COLON,
		"Expect ';' after value.",
	)
	if err != nil {
		return nil
	}

	return &PrintStmt{Expression: expr}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	if p.panicMode {
		return nil
	}

	_, err := p.consume(
		scanner.TOKEN_SEMI_COLON,
		"Expect ';' after expression.",
	)
	if err != nil {
		return nil
	}

	return &ExpressionStmt{Expression: expr}
}
