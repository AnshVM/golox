package Ast

import "github.com/AnshVM/golox/Tokens"

type Stmt interface {
	stmt()
}

type ExpressionStmt struct {
	Expression Expr
}

func (expr ExpressionStmt) stmt() {}

type PrintStmt struct {
	Expression Expr
}

func (expr PrintStmt) stmt() {}

type VarStmt struct {
	Name        *Tokens.Token
	Initializer Expr
}

func (expr VarStmt) stmt() {}

type BlockStmt struct {
	Statements []Stmt
}

func (b BlockStmt) stmt() {}
