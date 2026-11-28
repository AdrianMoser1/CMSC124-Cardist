package parser

import "cardist/scanner"

// forStatement parses after `for`:
//
//	"(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement
//
// There is no ForStmt node. It is desugared right here into nodes that already
// exist, so the evaluator needs no `for` case:
//
// The increment stays on the desugared while so continue runs it and break skips it.
func (p *Parser) forStatement() Stmt {
	if _, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'for'."); err != nil {
		return nil
	}

	// Clause 1: initializer. varDeclaration/expressionStatement eat their own ';'.
	var initializer Stmt
	if p.match(scanner.TOKEN_SEMI_COLON) {
		// empty clause: nothing to run first
	} else if p.match(scanner.TOKEN_VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}
	if p.panicMode {
		return nil
	}

	// Clause 2: condition. We eat the ';' ourselves.
	var condition Expr
	if !p.check(scanner.TOKEN_SEMI_COLON) {
		condition = p.expression()
		if p.panicMode {
			return nil
		}
	}
	if _, err := p.consume(scanner.TOKEN_SEMI_COLON, "Expect ';' after loop condition."); err != nil {
		return nil
	}

	// Clause 3: increment. Ends at ')' instead of ';'.
	var increment Expr
	if !p.check(scanner.TOKEN_RIGHT_PAREN) {
		increment = p.expression()
		if p.panicMode {
			return nil
		}
	}
	if _, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after for clauses."); err != nil {
		return nil
	}

	body := p.loopBody()
	if p.panicMode {
		return nil
	}

	if condition == nil {
		condition = &LiteralExpr{Value: true} // for (;;) loops forever
	}
	body = &WhileStmt{Condition: condition, Body: body, Increment: increment}
	if initializer != nil {
		// Outer block: keeps the loop variable from leaking past the loop.
		body = &BlockStmt{Statements: []Stmt{initializer, body}}
	}
	return body
}
