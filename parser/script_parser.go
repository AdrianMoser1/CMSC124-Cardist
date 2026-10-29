package parser

import "cardist/scanner"

// ParseScript parses statements for plain file execution.
// Parse() remains the expression-mode entry point for earlier labs.
func (p *Parser) ParseScript() *Script {
	script := &Script{}

	for !p.isAtEnd() {
		stmt := p.statement()

		// Week 1: reject the script at the first syntax error.
		// The caller checks ErrorFound() before executing anything.
		if p.hadError {
			break
		}

		if stmt != nil {
			script.Statements = append(script.Statements, stmt)
		}
	}

	return script
}

func (p *Parser) statement() Stmt {
	if p.match(scanner.TOKEN_PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	if p.hadError {
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
	if p.hadError {
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
