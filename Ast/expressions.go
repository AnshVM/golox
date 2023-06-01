package Ast

import (
	"github.com/AnshVM/golox/Tokens"
)

type Token = Tokens.Token

type Expr interface {
	isExpr()
}

type ConditionalExpr struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

func (c ConditionalExpr) isExpr() {}

type BinaryExpr struct {
	Left     Expr
	Operator *Token
	Right    Expr
}

func (b BinaryExpr) isExpr() {}

type GroupingExpr struct {
	Expression Expr
}

func (g GroupingExpr) isExpr() {}

type LiteralExpr struct {
	Value any
}

func (l LiteralExpr) isExpr() {}

type UnaryExpr struct {
	Operator *Token
	Right    Expr
}

func (u UnaryExpr) isExpr() {}

type VariableExpr struct {
	Name *Tokens.Token
}

func (v VariableExpr) isExpr() {}

type AssignExpr struct {
	Name  *Tokens.Token
	Value Expr
}

func (a *AssignExpr) isExpr() {}

type LogicalExpr struct {
	Left     Expr
	Operator *Tokens.Token
	Right    Expr
}

func (l *LogicalExpr) isExpr() {}

type Call struct {
	Callee    Expr
	Paren     *Tokens.Token
	Arguments []Expr
}

func (c *Call) isExpr() {}
