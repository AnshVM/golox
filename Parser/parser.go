package Parser

import (
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

type Parser struct {
	tokens     []*Token
	current    uint
	parseError error
}

func NewParser(tokens []*Tokens.Token) *Parser {
	return &Parser{
		tokens:     tokens,
		current:    0,
		parseError: nil,
	}
}

func (p *Parser) Parse() Expr {
	expr := p.comma()
	if p.parseError != nil {
		p.parseError = nil
		return nil
	}
	return expr
}

// Recursive decent parser
// Lower precedence in taken first

func (p *Parser) comma() Expr {
	expr := p.expression()

	for p.match(Tokens.COMMA) {
		operator := p.previous()
		right := p.expression()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparision()

	for p.match(Tokens.EQUAL_EQUAL, Tokens.BANG_EQUAL) {
		operator := p.previous()
		right := p.comparision()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparision() Expr {
	expr := p.term()

	for p.match(Tokens.GREATER, Tokens.GREATER_EQUAL, Tokens.LESS, Tokens.LESS_EQAUL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(Tokens.MINUS, Tokens.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(Tokens.SLASH, Tokens.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	var expr Expr
	for p.match(Tokens.MINUS, Tokens.BANG) {
		prefix := p.previous()
		right := p.unary()
		expr = &Unary{Operator: prefix, Right: right}
		return expr
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(Tokens.NUMBER, Tokens.STRING) {
		return &Literal{Value: p.previous().Literal}
	}

	if p.match(Tokens.TRUE) {
		return &Literal{Value: true}
	}
	if p.match(Tokens.FALSE) {
		return &Literal{Value: false}
	}
	if p.match(Tokens.NIL) {
		return &Literal{Value: nil}
	}

	if p.match(Tokens.LEFT_PAREN) {
		expr := p.expression()
		p.consume(Tokens.RIGHT_PAREN, "Expect ')' after expression")
		return &Grouping{Expression: expr}
	}
	p.parseError = Error.ErrParseError
	Error.ReportParseError(p.peek(), "Unexpected token")
	return nil
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == Tokens.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case Tokens.CLASS:
		case Tokens.FUN:
		case Tokens.VAR:
		case Tokens.FOR:
		case Tokens.IF:
		case Tokens.WHILE:
		case Tokens.PRINT:
		case Tokens.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(tokenType string, message string) {
	if p.check(tokenType) {
		p.advance()
		return
	}
	Error.ReportParseError(p.peek(), message)
	p.parseError = Error.ErrParseError
}

func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == Tokens.EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}
