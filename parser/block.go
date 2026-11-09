package parser

import "cardist/scanner"

// block → "{" declaration* "}"   (the '{' was already consumed by statement())
//
// Each declaration() recovers from its own errors, so one bad statement does
// not stop the rest of the block from being checked. At EOF the consume below
// reports "Expect '}' after block." (an unclosed block).
func (p *Parser) block() []Stmt {
	p.blockDepth++
	defer func() { p.blockDepth-- }()

	var statements []Stmt
	for !p.check(scanner.TOKEN_RIGHT_BRACE) && !p.isAtEnd() {
		if stmt := p.declaration(); stmt != nil {
			statements = append(statements, stmt)
		}
	}
	p.consume(scanner.TOKEN_RIGHT_BRACE, "Expect '}' after block.")
	return statements
}
