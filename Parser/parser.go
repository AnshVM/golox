package Parser

import (
	"fmt"

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

func (p *Parser) Parse() []Stmt {
	statements := []Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	if p.parseError != nil {
		p.parseError = nil
		return nil
	}
	return statements
}

// Recursive decent parser
// Lower precedence in taken first

func (p *Parser) declaration() Stmt {
	defer func() {
		if p.parseError != nil {
			p.parseError = nil
			p.synchronize()
		}
	}()
	if p.match(Tokens.VAR) {
		varName := p.consume(Tokens.IDENTIFIER, "Expected variable name")
		if varName == nil {
			return nil
		}
		if p.match(Tokens.EQUAL) {
			expr := p.expression()
			p.consume(Tokens.SEMICOLON, "Expect ';' after declaration")
			return Var{Name: varName, Initializer: expr}
		}
		p.consume(Tokens.SEMICOLON, "Expect ';' after declaration")
		return Var{Name: varName, Initializer: nil}
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	if p.match(Tokens.PRINT) {
		expr := p.expression()
		p.consume(Tokens.SEMICOLON, "Expect ';' after expression")
		return Print{Expression: expr}
	} else if p.match(Tokens.LEFT_BRACE) {
		statements := []Stmt{}
		for !p.check(Tokens.RIGHT_BRACE) && !p.isAtEnd() {
			statements = append(statements, p.declaration())
		}
		p.consume(Tokens.RIGHT_BRACE, "Expect '}' after block.")
		return Block{Statements: statements}
	} else {
		expr := p.expression()
		p.consume(Tokens.SEMICOLON, "Expect ';' after expression")
		return Expression{Expression: expr}
	}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | equality
func (p *Parser) assignment() Expr {
	expr := p.ternary()
	if p.match(Tokens.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if varExpr, ok := expr.(*Variable); ok {
			return &Assign{Name: varExpr.Name, Value: value}
		}
		Error.ReportParseError(equals, "Invalid assignment target")
	}
	return expr
}

// func (p *Parser) comma() Expr {

// 	if p.match(Tokens.COMMA) {
// 		p.missingExpressionBefore(",")
// 	}

// 	expr := p.ternary()

// 	for p.match(Tokens.COMMA) {
// 		operator := p.previous()
// 		right := p.ternary()
// 		expr = &Binary{Left: expr, Operator: operator, Right: right}
// 	}

// 	return expr
// }

func (p *Parser) ternary() Expr {
	expr := p.equality()

	if p.match(Tokens.QUESTION_MARK) {
		thenExpr := p.expression()
		p.consume(Tokens.COLON, "Expected ':' when using ternary operator '?'")
		elseExpr := p.expression()
		expr = &Conditional{Condition: expr, Then: thenExpr, Else: elseExpr}
	}
	return expr
}

func (p *Parser) equality() Expr {

	if p.match(Tokens.EQUAL_EQUAL) {
		p.missingExpressionBefore("==")
	}

	if p.match(Tokens.BANG_EQUAL) {
		p.missingExpressionBefore("!=")
	}
	expr := p.comparision()

	for p.match(Tokens.EQUAL_EQUAL, Tokens.BANG_EQUAL) {
		operator := p.previous()
		right := p.comparision()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) comparision() Expr {

	switch true {
	case p.match(Tokens.GREATER):
		p.missingExpressionBefore(">")
		break
	case p.match(Tokens.GREATER_EQUAL):
		p.missingExpressionBefore(">=")
		break
	case p.match(Tokens.LESS):
		p.missingExpressionBefore("<")
		break
	case p.match(Tokens.LESS_EQAUL):
		p.missingExpressionBefore("<=")
		break
	default:
		break
	}

	expr := p.term()

	for p.match(Tokens.GREATER, Tokens.GREATER_EQUAL, Tokens.LESS, Tokens.LESS_EQAUL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) term() Expr {

	switch true {
	case p.match(Tokens.MINUS):
		p.missingExpressionBefore("-")
		break
	case p.match(Tokens.PLUS):
		p.missingExpressionBefore("+")
		break
	default:
		break
	}

	expr := p.factor()

	for p.match(Tokens.MINUS, Tokens.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() Expr {
	switch true {
	case p.match(Tokens.SLASH):
		p.missingExpressionBefore("/")
		break
	case p.match(Tokens.STAR):
		p.missingExpressionBefore("*")
		break
	default:
		break
	}
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

	if p.match(Tokens.IDENTIFIER) {
		return &Variable{Name: p.previous()}
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

// Similar to `check()`, but accepts a list of token types
func (p *Parser) match(tokenTypes ...string) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// Consumes the current token and returns it
func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

// similar to `check()`, but throws an error with `message`
func (p *Parser) consume(tokenType string, message string) *Tokens.Token {
	if p.check(tokenType) {
		p.advance()
		return p.previous()
	}
	Error.ReportParseError(p.peek(), message)
	p.parseError = Error.ErrParseError
	return nil
}

// checks the type of the current token, consumes the token if types match
func (p *Parser) check(tokenType string) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == Tokens.EOF
}

// returns currrent token without consuming it
func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) missingExpressionBefore(operator string) {
	Error.ReportParseError(p.previous(), fmt.Sprintf("Missing expression before '%s'", operator))
	p.parseError = Error.ErrParseError
}
