package scanner

import (
	"fmt"
	"strconv"
	"unicode"
)

type TokenType int

const (
	// Single-character tokens.
	T_LEFT_PAREN TokenType = iota
	T_RIGHT_PAREN
	T_LEFT_BRACE
	T_RIGHT_BRACE
	T_COMMA
	T_DOT
	T_MINUS
	T_PLUS
	T_SEMICOLON
	T_SLASH
	T_STAR

	// One or two character tokens.
	T_BANG
	T_BANG_EQUAL
	T_EQUAL
	T_EQUAL_EQUAL
	T_GREATER
	T_GREATER_EQUAL
	T_LESS
	T_LESS_EQUAL

	// Literals.
	T_IDENTIFIER
	T_STRING
	T_NUMBER

	// Keywords.
	T_AND
	T_CLASS
	T_ELSE
	T_FALSE
	T_FUN
	T_FOR
	T_IF
	T_NIL
	T_OR
	T_PRINT
	T_RETURN
	T_SUPER
	T_THIS
	T_TRUE
	T_VAR
	T_WHILE

	T_EOF
)

//go:generate stringer -type=TokenType

// Token is token, yes
type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t *Token) Type() TokenType      { return t.tokenType }
func (t *Token) Lexeme() string       { return t.lexeme }
func (t *Token) Literal() interface{} { return t.literal }

func (t *Token) String() string {
	return fmt.Sprintf("Token(%s, %s)", t.tokenType, t.lexeme)
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
		} else if s.match('*') {
			lvl := 1
			for {
				if s.match('/') && s.match('*') {
					lvl++
				} else if s.match('*') && s.match('/') {
					lvl--
					if lvl == 0 {
						break
					}
				} else {
					s.advance()
				}
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
		if unicode.IsDigit(s.last()) {
			s.number()
			break
		}
		if unicode.IsLetter(s.last()) {
			s.identifier()
			break
		}
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

func (s *Scanner) number() {
	for unicode.IsDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		s.advance()
		for unicode.IsDigit(s.advance()) {
		}
	}
	val, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	s.addTokenWithLiteral(T_NUMBER, val)
}

func (s *Scanner) identifier() {
	for unicode.IsLetter(s.peek()) || unicode.IsDigit(s.peek()) {
		s.advance()
	}
	ident := s.source[s.start:s.current]
	switch string(ident) {
	case "and":
		s.addToken(T_AND)
	case "class":
		s.addToken(T_CLASS)
	case "else":
		s.addToken(T_ELSE)
	case "false":
		s.addToken(T_FALSE)
	case "for":
		s.addToken(T_FOR)
	case "fun":
		s.addToken(T_FUN)
	case "if":
		s.addToken(T_IF)
	case "nil":
		s.addToken(T_NIL)
	case "or":
		s.addToken(T_OR)
	case "print":
		s.addToken(T_PRINT)
	case "return":
		s.addToken(T_RETURN)
	case "super":
		s.addToken(T_SUPER)
	case "this":
		s.addToken(T_IF)
	case "true":
		s.addToken(T_TRUE)
	case "var":
		s.addToken(T_VAR)
	case "while":
		s.addToken(T_WHILE)
	default:
		s.addToken(T_IDENTIFIER)
	}
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
func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return rune(0)
	}
	return s.source[s.current+1]
}
func (s *Scanner) match(to rune) bool {
	if s.peek() == to {
		s.advance()
		return true
	}
	return false
}
