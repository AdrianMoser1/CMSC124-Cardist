package scanner

import (
	"fmt"
	"os"
)

type TokenType int

const (
	//single
	TOKEN_LEFT_PAREN TokenType = iota //iota so each const var gets assigned a sequential number
	TOKEN_RIGHT_PAREN
	TOKEN_LEFT_BRACE  // {
	TOKEN_RIGHT_BRACE // }
	TOKEN_PLUS        // +
	TOKEN_MINUS       // -
	TOKEN_STAR        // *
	TOKEN_EQUAL       // =
	TOKEN_GREATER     // >
	TOKEN_LESSER      // <
	TOKEN_COMMA       // ,
	TOKEN_COLON       // :
	TOKEN_SEMI_COLON  // ;


	//double
	TOKEN_BANG //!
	TOKEN_BANG_EQUAL //!=
	TOKEN_EQUAL_EQUAL //==
	TOKEN_GREATER_EQUAL //>=
	TOKEN_LESSER_EQUAL //<=
	TOKEN_SLASH // / 

	TOKEN_EOF
)

func (t TokenType) String() string {
	switch t {
	case TOKEN_LEFT_PAREN:
		return "LEFT_PAREN"
	case TOKEN_RIGHT_PAREN:
		return "RIGHT_PAREN"
	case TOKEN_LEFT_BRACE:
		return "LEFT_BRACE"
	case TOKEN_RIGHT_BRACE:
		return "RIGHT_BRACE"
	case TOKEN_PLUS:
		return "PLUS"
	case TOKEN_MINUS:
		return "MINUS"
	case TOKEN_STAR:
		return "STAR"
	case TOKEN_EQUAL:
		return "EQUAL"
	case TOKEN_GREATER:
		return "GREATER"
	case TOKEN_LESSER:
		return "LESSER"
	case TOKEN_COMMA:
		return "COMMA"
	case TOKEN_COLON:
		return "COLON"
	case TOKEN_SEMI_COLON:
		return "SEMI_COLON"
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_BANG:
		return "BANG"
	case TOKEN_BANG_EQUAL:
		return "BANG_EQUAL"
	case TOKEN_EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case TOKEN_GREATER_EQUAL:
		return "GREATER_EQUAL"
	case TOKEN_LESSER_EQUAL:
		return "LESSER_EQUAL"
	case TOKEN_SLASH:
		return "SLASH"
	default:
		return "UNKNOWN"
	}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}
type Scanner struct {
	source   string
	tokens   []Token //tokens collected
	start    int
	current  int
	line     int
	hadError bool
}

func (s *Scanner) ErrorFound() bool { // Reports whether any error was recorded during scanning.
	return s.hadError // Used this to decide the exit code (0 vs 65).
}

/*
	reportError prints a diagnostic to stderr and marks the scan as failed.

It does not stop scanning, so later errors in the same file can still be reported in one pass.
*/
func (s *Scanner) reportError(line int, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)
	s.hadError = true
}

func (t Token) String() string {
	literalStr := "null"
	if t.Literal != nil { //will return as it is if it has value
		literalStr = fmt.Sprintf("%v", t.Literal)
	}
	return fmt.Sprintf("Token(type=%s, lexeme=%s, literal=%s, line=%d)", t.Type, t.Lexeme, literalStr, t.Line)
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool { //will return true if current reaches end of source string per line
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte { //stepping forward by 1 whilst returning character at current
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) addToken(t TokenType) { //creates token type and appends to tokens or # of tokens
	s.addTokenWithLiteral(t, nil)
}

func (s *Scanner) match(expected byte) bool { //for consumation of characters only if matching
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() byte { //to check current character without advancing current
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNex() byte {
	if s.current+1 >= len(s.source){
		return 0
	}
	return s.source[s.current+1]
}


func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    s.line,
	})
}

// Main scanner feature
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{ //appending eof
		Type:    TOKEN_EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	//single-char
	case '(':
		s.addToken(TOKEN_LEFT_PAREN)
	case ')':
		s.addToken(TOKEN_RIGHT_PAREN)
	case '{':
		s.addToken(TOKEN_LEFT_BRACE)
	case '}':
		s.addToken(TOKEN_RIGHT_BRACE)
	case '+':
		s.addToken(TOKEN_PLUS)
	case '-':
		s.addToken(TOKEN_MINUS)
	case '*':
		s.addToken(TOKEN_STAR)
	case ',':
		s.addToken(TOKEN_COMMA)
	case ':':
		s.addToken(TOKEN_COLON)
	case ';':
		s.addToken(TOKEN_SEMI_COLON)

	//double-char w/h match
	case '!':
		if s.match('=') {
			s.addToken(TOKEN_BANG_EQUAL)
		} else {
			s.addToken(TOKEN_BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(TOKEN_EQUAL_EQUAL)
		} else {
			s.addToken(TOKEN_EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(TOKEN_LESSER_EQUAL)
		} else {
			s.addToken(TOKEN_LESSER)
		}
	case '>':
		if s.match('=') {
			s.addToken(TOKEN_GREATER_EQUAL)
		} else {
			s.addToken(TOKEN_GREATER)
		}

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			// use peek() so we DON'T consume the '\n' character here;
			// the next iteration of scanToken() will handle '\n' and increment s.line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TOKEN_SLASH)
		}	
	
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++ //increment line count for error reporting
	default:
		s.reportError(s.line, fmt.Sprintf("Unexpected character: %c", c)) //report error but keep scanning, unincluded
	}
}
