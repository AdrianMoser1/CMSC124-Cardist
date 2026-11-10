package parser

import "cardist/scanner"

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (*IfStmt) stmtNode() {}

// ifStatement parses after `if`: "(" expression ")" statement ["else" statement].
// A chained `elseif` is parsed as an `else` branch containing a nested `if`.
func (p *Parser) ifStatement() Stmt {
	if _, err := p.consume(scanner.TOKEN_LEFT_PAREN, "Expect '(' after 'if'."); err != nil {
		return nil
	}
	condition := p.expression()
	if p.panicMode {
		return nil
	}
	if _, err := p.consume(scanner.TOKEN_RIGHT_PAREN, "Expect ')' after condition."); err != nil {
		return nil
	}

	thenBranch := p.statement()
	if p.panicMode {
		return nil
	}

	var elseBranch Stmt
	if p.match(scanner.TOKEN_ELSE) {
		if p.match(scanner.TOKEN_IF) {
			elseBranch = p.ifStatement()
		} else {
			elseBranch = p.statement()
		}
		if p.panicMode {
			return nil
		}
	} else if p.match(scanner.TOKEN_ELSEIF) {
		elseBranch = p.ifStatement()
		if p.panicMode {
			return nil
		}
	}

	return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}
