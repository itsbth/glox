package scanner

import "fmt"

type TokenType int

const (
	// Single-character tokens.
	T_LEFT_PAREN  TokenType = iota
	T_RIGHT_PAREN           = iota
	T_LEFT_BRACE            = iota
	T_RIGHT_BRACE           = iota
	T_COMMA                 = iota
	T_DOT                   = iota
	T_MINUS                 = iota
	T_PLUS                  = iota
	T_SEMICOLON             = iota
	T_SLASH                 = iota
	T_STAR                  = iota

	// One or two character tokens.
	T_BANG          = iota
	T_BANG_EQUAL    = iota
	T_EQUAL         = iota
	T_EQUAL_EQUAL   = iota
	T_GREATER       = iota
	T_GREATER_EQUAL = iota
	T_LESS          = iota
	T_LESS_EQUAL    = iota

	// Literals.
	T_IDENTIFIER = iota
	T_STRING     = iota
	T_NUMBER     = iota

	// Keywords.
	T_AND    = iota
	T_CLASS  = iota
	T_ELSE   = iota
	T_FALSE  = iota
	T_FUN    = iota
	T_FOR    = iota
	T_IF     = iota
	T_NIL    = iota
	T_OR     = iota
	T_PRINT  = iota
	T_RETURN = iota
	T_SUPER  = iota
	T_THIS   = iota
	T_TRUE   = iota
	T_VAR    = iota
	T_WHILE  = iota

	T_EOF = iota
)

func (t TokenType) String() string {
	switch t {
	case T_LEFT_PAREN:
		return "("
	default:
		return "tbd"
	}
}

// Token is token, yes
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

// UnexpectedTokenError is error during scanning, see
type UnexpectedTokenError struct {
	expected  []rune
	found     rune
	pos, line int
}

func (err UnexpectedTokenError) Error() string {
	if len(err.expected) == 1 {
		return fmt.Sprintf("found %c while looking for %c", err.found, err.expected)
	}
	return "tbd"
}

func unexpectedToken(found rune, expected ...rune) UnexpectedTokenError {
	return UnexpectedTokenError{
		found:    found,
		expected: expected,
		pos:      0,
		line:     0,
	}
}

// Scanner is scanner, yes
type Scanner struct {
	source         []rune
	current, start int
	line           int
	Tokens         chan Token
	Errors         chan error
}

// NewScanner creates a new scanner
func NewScanner(source string) *Scanner {
	return &Scanner{
		source:  []rune(source),
		current: 0,
		start:   0,
		line:    1,
		Tokens:  make(chan Token),
		Errors:  make(chan error),
	}
}

// Scan scans
func (s *Scanner) Scan() {
	for !s.eof() {
		s.start = s.current
		s.scanToken()
	}
	close(s.Tokens)
	close(s.Errors)
}

func (s *Scanner) scanToken() {
	switch s.advance() {
	case '(':
		s.addToken(T_LEFT_PAREN)
	case ')':
		s.addToken(T_RIGHT_PAREN)
	case '{':
		s.addToken(T_LEFT_BRACE)
	case '}':
		s.addToken(T_RIGHT_BRACE)
	case ',':
		s.addToken(T_COMMA)
	case '.':
		s.addToken(T_DOT)
	case '+':
		s.addToken(T_PLUS)
	case '-':
		s.addToken(T_MINUS)
	case '*':
		s.addToken(T_STAR)
	case ';':
		s.addToken(T_SEMICOLON)
	case '!':
		s.addTokenIf('=', T_BANG_EQUAL, T_BANG)
	case '=':
		s.addTokenIf('=', T_EQUAL_EQUAL, T_EQUAL)
	case '<':
		s.addTokenIf('=', T_LESS_EQUAL, T_LESS)
	case '>':
		s.addTokenIf('=', T_GREATER_EQUAL, T_GREATER)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.eof() {
				s.advance()
			}
		} else {
			s.addToken(T_SLASH)
		}
	case ' ':
	case '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		s.Errors <- fmt.Errorf("unknown character %c", s.last())
	}
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.eof() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}
	if s.eof() {
		s.Errors <- unexpectedToken('a', '"')
		return
	}
	s.advance()
	val := string(s.source[s.start+1 : s.current-1])
	s.addTokenWithLiteral(T_STRING, val)
}

func (s *Scanner) addToken(token TokenType) {
	s.addTokenWithLiteral(token, nil)
}

func (s *Scanner) addTokenWithLiteral(token TokenType, literal interface{}) {
	s.Tokens <- Token{tokenType: token, lexeme: string(s.source[s.start:s.current]), line: s.line, literal: literal}
}

func (s *Scanner) addTokenIf(next rune, two TokenType, one TokenType) {
	if s.match(next) {
		s.addToken(two)
	} else {
		s.addToken(one)
	}
}

func (s *Scanner) eof() bool     { return s.current >= len(s.source) }
func (s *Scanner) last() rune    { return s.source[s.current-1] }
func (s *Scanner) advance() rune { s.current++; return s.source[s.current-1] }
func (s *Scanner) peek() rune {
	if s.eof() {
		return rune(0)
	}
	return s.source[s.current]
}
func (s *Scanner) match(to rune) bool {
	if s.peek() == to {
		s.advance()
		return true
	}
	return false
}
