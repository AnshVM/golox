package Parser

import (
	"github.com/AnshVM/golox/Tokens"
)

type Token = Tokens.Token

type Expr interface {
	isExpr()
}

type Conditional struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

func (c Conditional) isExpr() {}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (b Binary) isExpr() {}

type Grouping struct {
	Expression Expr
}

func (g Grouping) isExpr() {}

type Literal struct {
	Value any
}

func (l Literal) isExpr() {}

type Unary struct {
	Operator *Token
	Right    Expr
}

func (u Unary) isExpr() {}

type Variable struct {
	Name *Tokens.Token
}

func (v Variable) isExpr() {}

type Assign struct {
	Name  *Tokens.Token
	Value Expr
}

func (a *Assign) isExpr() {}

type Stmt interface {
	stmt()
}

type Expression struct {
	Stmt
	Expression Expr
}

func (expr Expression) stmt() {}

type Print struct {
	Stmt
	Expression Expr
}

func (expr Print) stmt() {}

type Var struct {
	Stmt
	Name        *Tokens.Token
	Initializer Expr
}

func (expr Var) stmt() {}

type Block struct {
	Statements []Stmt
}

func (b Block) stmt() {}
