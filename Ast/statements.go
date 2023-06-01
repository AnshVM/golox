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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (ifs IfStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func (while WhileStmt) stmt() {}

type Function struct {
	Name   *Tokens.Token
	Params []*Tokens.Token
	Body   []Stmt
}

func (f Function) stmt() {}

type Return struct {
	Keyword *Tokens.Token
	Value   Expr
}

func (r Return) stmt() {}
