package Parser

import (
	"github.com/AnshVM/golox/Tokens"
)

type Token = Tokens.Token

type Expr interface {
	Print() string
	PrintRPN() string
}

type Binary struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value any
}

type Unary struct {
	Operator *Token
	Right    Expr
}
