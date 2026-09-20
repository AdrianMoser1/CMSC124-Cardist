package scanner

import (
	"fmt"
	"os"
	"strconv"
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
	TOKEN_BANG          //!
	TOKEN_BANG_EQUAL    //!=
	TOKEN_EQUAL_EQUAL   //==
	TOKEN_GREATER_EQUAL //>=
	TOKEN_LESSER_EQUAL  //<=
	TOKEN_SLASH         // /

	//literals
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER

	//literals-general
	TOKEN_VAR
	TOKEN_IF
	TOKEN_ELSE
	TOKEN_ELSEIF
	TOKEN_SWITCH
	TOKEN_WHILE
	TOKEN_TRUE
	TOKEN_FALSE
	TOKEN_NIL
	TOKEN_AND
	TOKEN_OR
	TOKEN_PRINT
	TOKEN_RETURN
	TOKEN_FOR
	TOKEN_FUNC
	TOKEN_BREAK
	TOKEN_CONTINUE

	//keywords-cardist domain
	TOKEN_CARD
	TOKEN_ENERGY
	TOKEN_COST
	TOKEN_ENEMY
	TOKEN_INTENT
	TOKEN_ARTIFACT
	TOKEN_ELIXIR
	TOKEN_DEAL
	TOKEN_BLOCK
	TOKEN_APPLY
	TOKEN_DRAW
	TOKEN_BANISH
	TOKEN_TURN
	TOKEN_PLAYER
	TOKEN_EFFECT

	//end of file

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
	case TOKEN_IDENTIFIER:
		return "IDENTIFIER"
	case TOKEN_STRING:
		return "STRING"
	case TOKEN_NUMBER:
		return "NUMBER"
	case TOKEN_VAR:
		return "VAR"
	case TOKEN_IF:
		return "IF"
	case TOKEN_ELSE:
		return "ELSE"
	case TOKEN_ELSEIF:
		return "ELSEIF"
	case TOKEN_SWITCH:
		return "SWITCH"
	case TOKEN_PRINT:
		return "PRINT"
	case TOKEN_RETURN:
		return "RETURN"
	case TOKEN_FOR:
		return "FOR"
	case TOKEN_FUNC:
		return "FUNC"
	case TOKEN_BREAK:
		return "BREAK"
	case TOKEN_CONTINUE:
		return "CONTINUE"
	case TOKEN_WHILE:
		return "WHILE"
	case TOKEN_TRUE:
		return "TRUE"
	case TOKEN_FALSE:
		return "FALSE"
	case TOKEN_NIL:
		return "NIL"
	case TOKEN_AND:
		return "AND"
	case TOKEN_OR:
		return "OR"
	case TOKEN_CARD:
		return "CARD"
	case TOKEN_ENERGY:
		return "ENERGY"
	case TOKEN_COST:
		return "COST"
	case TOKEN_ENEMY:
		return "ENEMY"
	case TOKEN_INTENT:
		return "INTENT"
	case TOKEN_ARTIFACT:
		return "ARTIFACT"
	case TOKEN_ELIXIR:
		return "ELIXIR"
	case TOKEN_DEAL:
		return "DEAL"
	case TOKEN_BLOCK:
		return "BLOCK"
	case TOKEN_APPLY:
		return "APPLY"
	case TOKEN_DRAW:
		return "DRAW"
	case TOKEN_BANISH:
		return "BANISH"
	case TOKEN_TURN:
		return "TURN"
	case TOKEN_PLAYER:
		return "PLAYER"
	case TOKEN_EFFECT:
		return "EFFECT"
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

func (s *Scanner) peekNext() byte { //to check next character without advancing current
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func isAlpha(c byte) bool { //checks if character is an alphabet
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isDigit(c byte) bool { //checks if character is digit
	return c >= '0' && c <= '9'
}

func isAlphaNumeric(c byte) bool { //checks if character is alphanumeric
	return isAlpha(c) || isDigit(c)
}

// Maps every reserved word to its token type
// Anything that starts with a letter but isn't in this table is a considered as an IDENTIFIER.
var keywords = map[string]TokenType{
	"print": TOKEN_PRINT, "return": TOKEN_RETURN, "for": TOKEN_FOR,
	"func": TOKEN_FUNC, "break": TOKEN_BREAK, "continue": TOKEN_CONTINUE,
	"elseif": TOKEN_ELSEIF, "switch": TOKEN_SWITCH, "effect": TOKEN_EFFECT,
	"var": TOKEN_VAR, "if": TOKEN_IF, "else": TOKEN_ELSE, "while": TOKEN_WHILE,
	"true": TOKEN_TRUE, "false": TOKEN_FALSE, "nil": TOKEN_NIL,
	"and": TOKEN_AND, "or": TOKEN_OR,
	"card": TOKEN_CARD, "energy": TOKEN_ENERGY, "cost": TOKEN_COST,
	"enemy": TOKEN_ENEMY, "intent": TOKEN_INTENT, "artifact": TOKEN_ARTIFACT,
	"elixir": TOKEN_ELIXIR, "deal": TOKEN_DEAL, "block": TOKEN_BLOCK,
	"apply": TOKEN_APPLY, "draw": TOKEN_DRAW, "banish": TOKEN_BANISH,
	"turn": TOKEN_TURN, "player": TOKEN_PLAYER,
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

// identifier consumes the WHOLE run of letters/digits/underscores first (by using peek())
// then compares the finished string against the keyword table
// Checking keywords one character at a time instead would break
func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current] //final string of the identifier, to be checked against keywords table
	if t, ok := keywords[text]; ok {    //if the text is found in the keywords map, then it's a keyword
		s.addToken(t)
	} else {
		s.addToken(TOKEN_IDENTIFIER)
	}
}

// number handles both integers and decimals. A "." is only part of the
// number if a digit follows it — otherwise it's its own token later,
// so "3.toString" and "3.14" both scan sensibly.
func (s *Scanner) number() {
	for isDigit(s.peek()) { //consume all digits
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) { //next number is a decimal, so consume the "." and the following digits
		s.advance() // consume the "."
		for isDigit(s.peek()) {
			s.advance()
		}
	}
	text := s.source[s.start:s.current]
	value, err := strconv.ParseFloat(text, 64) //conversion to float64
	if err != nil {
		s.reportError(s.line, fmt.Sprintf("Invalid number literal: %s", text))
		return
	}
	s.addTokenWithLiteral(TOKEN_NUMBER, value)
}

// string consumes up to the closing quote. Per README, Cardist strings
// may not span lines — a newline before the closing quote is reported as an error
// The literal value of the string excludes the bordering quotes.
func (s *Scanner) string() {
	startLine := s.line

	for s.peek() != '"' && !s.isAtEnd() { //checks for the closing quote and if it has reached the end of the string
		if s.peek() == '\n' { // error if newline is found before closing quote
			s.reportError(startLine, "Unterminated string.")
			return
		}
		s.advance()
	}

	if s.isAtEnd() { //if it has reached the end of the string without finding a closing quote
		s.reportError(startLine, "Unterminated string.")
		return
	}

	s.advance() // consume the closing "

	// Literal excludes the bordering quotes; while lexeme keeps them.
	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(TOKEN_STRING, value)
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
	case '"':
		s.string()
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++ //increment line count for error reporting
	default: //check if the character is a digit or an alphabetical character or if its an unexpected character
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.reportError(s.line, fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}
