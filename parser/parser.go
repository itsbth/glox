package parser

import (
	"github.com/itsbth/glox/ast"
	"github.com/itsbth/glox/scanner"
)

type Parser struct {
	tokens  []scanner.Token
	current int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ast.Expr { return p.expression() }

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	// lhs := p.comparison()
	return nil
}

func (p *Parser) match(typ ...scanner.TokenType) bool {
	top := p.peek()
	for _, tok := range typ {
		if top.Type() == tok {
			p.advance()
			return true
		}
	}
	return false
}
func (p *Parser) peek() scanner.Token { return p.tokens[p.current] }
func (p *Parser) advance()            { p.current++ }
