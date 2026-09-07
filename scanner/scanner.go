package scanner

import "fmt" 

type TokenType int


// 
const ( 
	TOKEN_LEFT_PAREN 	TokenType = iota //iota so each const var gets assigned a sequential number
	TOKEN_RIGHT_PAREN 
	TOKEN_LEFT_BRACE // {
	TOKEN_RIGHT_BRACE // }
	TOKEN_PLUS // +
	TOKEN_MINUS // -
	TOKEN_STAR // * 
	TOKEN_EQUAL // =
	TOKEN_GREATER // >
	TOKEN_LESSER // < 
	TOKEN_COMMA // ,
	TOKEN_COLON // :
	TOKEN_SEMI_COLON // ;
	TOKEN_EOF
)

func (t TokenType) String() string{
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
	default: 
		return "UNKNOWN"
	}
}
type Token struct {
	Type	TokenType
	Lexeme	string
	Literal	any
	Line	int
}
type Scanner struct{
	source	string
	tokens	[]Token //tokens collected
	start	int
	current	int
	line	int
}

func (t Token) String() string{
	literalStr := "null"
	if t.Literal != nil { //will return as it is if it has value
		literalStr = fmt.Sprintf("%v", t.Literal)
	}
	return fmt.Sprintf("Token(type=%s, lexeme=%s, literal=%s, line=%d)", t.Type, t.Lexeme, literalStr, t.Line)	
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source:		source,
		tokens:		[]Token{},
		start:		0,
		current:	0,
		line:		1,
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

func (s *Scanner) addTokenWithLiteral(t TokenType, literal any) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		Type:		t,
		Lexeme:		lexeme,
		Literal:	literal, 
		Line:		s.line,
	})
}

//Main scanner feature
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, Token{ //appending eof 
		Type:		TOKEN_EOF,
		Lexeme:		"",
		Literal:	nil,
		Line:		s.line,
	})
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
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
	case '=':
		s.addToken(TOKEN_EQUAL)
	case '>':
		s.addToken(TOKEN_GREATER)
	case '<':
		s.addToken(TOKEN_LESSER)
	case ',':
		s.addToken(TOKEN_COMMA)
	case ':':
		s.addToken(TOKEN_COLON)
	case ';':
		s.addToken(TOKEN_SEMI_COLON)
	
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	default: 
		fmt.Printf("[line %d] Error: Unexpected character: %c/n", s.line, c)
	}
}