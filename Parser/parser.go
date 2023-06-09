package Parser

import (
	"fmt"

	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

type Token = Tokens.Token
type Stmt = Ast.Stmt
type Expr = Ast.Expr

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
	switch true {
	case p.match(Tokens.VAR):
		return p.varDecl()
	case p.match(Tokens.FUN):
		return p.funcDecl()
	default:
		return p.statement()
	}
}

func (p *Parser) funcDecl() Stmt {
	return p.function()
}

func (p *Parser) function() Stmt {
	if p.check(Tokens.IDENTIFIER) {
		return p.namedFunction(p.advance(), "function")
	} else {
		return p.expressionStmt()
	}
}

func (p *Parser) paramList(paren *Tokens.Token, kind string) []*Tokens.Token {
	var paramList []*Tokens.Token
	if p.peek().Type != Tokens.RIGHT_PAREN {
		paramList = p.params()
	}
	if len(paramList) >= 255 {
		Error.ReportParseError(paren, "Can't have more than 255 arguments")
		p.parseError = Error.ErrParseError
	}
	p.consume(Tokens.RIGHT_PAREN, fmt.Sprintf("Expect ')' after %s declaration", kind))
	return paramList
}

func (p *Parser) anonymousFunction(kind string) Expr {
	paren := p.consume(Tokens.LEFT_PAREN, fmt.Sprintf("Expect '(' after fun"))
	params := p.paramList(paren, "function")
	p.consume(Tokens.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", kind))
	stmts := p.block()
	return &Ast.AnonymousFuncion{Params: params, Body: stmts}
}

func (p *Parser) namedFunction(name *Token, kind string) Stmt {
	paren := p.consume(Tokens.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name", kind))
	params := p.paramList(paren, "function")
	p.consume(Tokens.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", kind))
	stmts := p.block()
	return &Ast.NamedFunction{Name: name, Params: params, Body: stmts}
}

func (p *Parser) params() []*Tokens.Token {
	paramList := []*Tokens.Token{}
	for {
		param := p.consume(Tokens.IDENTIFIER, "Expect parameter name.")
		paramList = append(paramList, param)
		if !p.match(Tokens.COMMA) {
			break
		}
	}
	return paramList
}

func (p *Parser) varDecl() Stmt {
	varName := p.consume(Tokens.IDENTIFIER, "Expected variable name")
	if varName == nil {
		return nil
	}
	if p.match(Tokens.EQUAL) {
		expr := p.expression()
		p.consume(Tokens.SEMICOLON, "Expect ';' after declaration")
		return &Ast.VarStmt{Name: varName, Initializer: expr}
	}
	p.consume(Tokens.SEMICOLON, "Expect ';' after declaration")
	return &Ast.VarStmt{Name: varName, Initializer: nil}
}

func (p *Parser) statement() Stmt {
	switch true {
	case p.match(Tokens.PRINT):
		return p.print()
	case p.match(Tokens.LEFT_BRACE):
		return &Ast.BlockStmt{Statements: p.block()}
	case p.match(Tokens.IF):
		return p.ifStmt()
	case p.match(Tokens.WHILE):
		return p.WhileStmt()
	case p.match(Tokens.FOR):
		return p.ForStmt()
	case p.match(Tokens.RETURN):
		return p.ReturnStmt()
	default:
		return p.expressionStmt()
	}
}

func (p *Parser) ReturnStmt() Stmt {
	keyword := p.previous()
	if p.match(Tokens.SEMICOLON) {
		return &Ast.Return{Keyword: keyword, Value: nil}
	}
	expr := p.expression()
	p.consume(Tokens.SEMICOLON, "Expect ';' after return.")
	return &Ast.Return{Keyword: keyword, Value: expr}
}

// desugarises to While loop
func (p *Parser) ForStmt() Stmt {
	p.consume(Tokens.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(Tokens.SEMICOLON) {
		initializer = nil
	} else if p.peek().Type == Tokens.VAR {
		initializer = p.declaration()
	} else {
		initializer = p.expressionStmt()
	}
	var condition Expr
	if !p.check(Tokens.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(Tokens.SEMICOLON, "Expect ';' after loop condition")

	var increment Expr
	if !p.check(Tokens.RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(Tokens.RIGHT_PAREN, "Expect ')' after for clauses")
	incrementStmt := &Ast.ExpressionStmt{Expression: increment}

	body := p.statement()

	if increment != nil {
		body = &Ast.BlockStmt{Statements: []Stmt{body, incrementStmt}}
	}
	if condition == nil {
		condition = &Ast.LiteralExpr{Value: true}
	}
	body = &Ast.WhileStmt{Condition: condition, Body: body}
	if initializer != nil {
		body = &Ast.BlockStmt{Statements: []Stmt{initializer, body}}
	}
	return body
}

func (p *Parser) WhileStmt() Stmt {
	p.consume(Tokens.LEFT_PAREN, "Expect '(' after 'while'.")
	expr := p.expression()
	p.consume(Tokens.RIGHT_PAREN, "Expect ')' after expression")
	body := p.statement()
	return &Ast.WhileStmt{Condition: expr, Body: body}
}

func (p *Parser) ifStmt() Stmt {
	p.consume(Tokens.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(Tokens.RIGHT_PAREN, "Expect ')' after expression.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(Tokens.ELSE) {
		elseBranch = p.statement()
	}
	return &Ast.IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) print() Stmt {
	expr := p.expression()
	p.consume(Tokens.SEMICOLON, "Expect ';' after expression")
	return &Ast.PrintStmt{Expression: expr}
}

func (p *Parser) block() []Stmt {
	statements := []Stmt{}
	for !p.check(Tokens.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(Tokens.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) expressionStmt() Stmt {
	expr := p.expression()
	p.consume(Tokens.SEMICOLON, "Expect ';' after expression")
	return &Ast.ExpressionStmt{Expression: expr}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

// assignment -> IDENTIFIER "=" assignment | equality
func (p *Parser) assignment() Expr {
	expr := p.funcExpr()
	if p.match(Tokens.EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if varExpr, ok := expr.(*Ast.VariableExpr); ok {
			return &Ast.AssignExpr{Name: varExpr.Name, Value: value}
		}
		Error.ReportParseError(equals, "Invalid assignment target")
	}
	return expr
}

func (p *Parser) funcExpr() Expr {
	if p.match(Tokens.FUN) {
		paren := p.consume(Tokens.LEFT_PAREN, "Expect '(' after fun")
		params := p.paramList(paren, "function")
		p.consume(Tokens.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body", "function"))
		stmts := p.block()
		return &Ast.AnonymousFuncion{Params: params, Body: stmts}
	}
	return p.logic_or()
}

func (p *Parser) logic_or() Expr {
	expr := p.logic_and()
	for p.match(Tokens.OR) {
		operator := p.previous()
		right := p.logic_and()
		expr = &Ast.LogicalExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) logic_and() Expr {
	expr := p.equality()
	for p.match(Tokens.AND) {
		operator := p.previous()
		right := p.equality()
		expr = &Ast.LogicalExpr{Left: expr, Operator: operator, Right: right}
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
		expr = &Ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
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
		expr = &Ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
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
		expr = &Ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
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
		expr = &Ast.BinaryExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() Expr {
	var expr Expr
	for p.match(Tokens.MINUS, Tokens.BANG) {
		prefix := p.previous()
		right := p.unary()
		expr = &Ast.UnaryExpr{Operator: prefix, Right: right}
		return expr
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for p.match(Tokens.LEFT_PAREN) {
		token := p.previous()
		if p.match(Tokens.RIGHT_PAREN) { //no args
			expr = &Ast.Call{Callee: expr, Arguments: []Expr{}, Paren: token}
			continue
		}
		args := []Expr{}
		for {
			arg := p.expression()
			args = append(args, arg)
			if !p.match(Tokens.COMMA) {
				break
			}
		}
		p.consume(Tokens.RIGHT_PAREN, "Expect ')' for function call.")

		if len(args) >= 255 {
			Error.ReportParseError(p.peek(), "Can't have more that 255 arguments")
		}
		expr = &Ast.Call{Callee: expr, Arguments: args, Paren: token}
	}
	return expr
}

func (p *Parser) primary() Expr {
	if p.match(Tokens.NUMBER, Tokens.STRING) {
		return &Ast.LiteralExpr{Value: p.previous().Literal}
	}

	if p.match(Tokens.TRUE) {
		return &Ast.LiteralExpr{Value: true}
	}
	if p.match(Tokens.FALSE) {
		return &Ast.LiteralExpr{Value: false}
	}
	if p.match(Tokens.NIL) {
		return &Ast.LiteralExpr{Value: nil}
	}

	if p.match(Tokens.IDENTIFIER) {
		return &Ast.VariableExpr{Name: p.previous()}
	}

	if p.match(Tokens.LEFT_PAREN) {
		expr := p.expression()
		p.consume(Tokens.RIGHT_PAREN, "Expect ')' after expression")
		return &Ast.GroupingExpr{Expression: expr}
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

// Similar to `check()`, but accepts a list of token types and consumes
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

// checks the type of the current token, does not consume
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
