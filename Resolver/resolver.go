package Resolver

import (
	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Interpreter"
	"github.com/AnshVM/golox/Tokens"
	"github.com/AnshVM/golox/Utils"
)

const (
	FUNCTION = iota
	NONE     = iota
)

type Resolver struct {
	interpreter     *Interpreter.Interpreter
	scopes          Utils.Stack[map[string]bool]
	currentFunction int
}

func NewResolver(interpreter *Interpreter.Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          Utils.NewStack[map[string]bool](),
		currentFunction: NONE,
	}
}

func (r *Resolver) Resolve(node any) {
	switch n := node.(type) {

	case []Ast.Stmt:
		for _, stmt := range n {
			r.Resolve(stmt)
		}
		break
	case *Ast.BlockStmt:
		r.beginScope()
		r.Resolve(n.Statements)
		r.endScope()
		break
	case *Ast.VarStmt:
		r.declare(n.Name)
		if n.Initializer != nil {
			r.Resolve(n.Initializer)
		}
		r.define(n.Name)
		break
	case *Ast.VariableExpr:
		scope, err := r.scopes.Peek()
		if err == nil {
			if declared, ok := scope[n.Name.Lexeme]; ok && !declared {
				Error.ReportParseError(n.Name, "Can't read local variable in its own initializer.")
			}
		}
		r.resolveLocal(n, n.Name)
		break
	case *Ast.AssignExpr:
		r.Resolve(n.Value)
		r.resolveLocal(n, n.Name)
		break

	case *Ast.NamedFunction:
		r.declare(n.Name)
		r.define(n.Name)
		r.resolveFunction(n)
		break

	case *Ast.AnonymousFuncion:
		enclosingFunction := r.currentFunction
		r.currentFunction = FUNCTION
		r.beginScope()
		for _, arg := range n.Params {
			r.declare(arg)
			r.define(arg)
		}
		r.Resolve(n.Body)
		r.endScope()
		r.currentFunction = enclosingFunction
		break

	case *Ast.ExpressionStmt:
		r.Resolve(n.Expression)
		break

	case *Ast.IfStmt:
		r.Resolve(n.Condition)
		r.Resolve(n.ThenBranch)
		if n.ElseBranch != nil {
			r.Resolve(n.ElseBranch)
		}
		break

	case *Ast.PrintStmt:
		r.Resolve(n.Expression)
		break

	case *Ast.Return:
		if r.currentFunction == NONE {
			Error.ReportParseError(n.Keyword, "Cannot return from top-level code.")
		}
		r.Resolve(n.Value)
		break

	case *Ast.WhileStmt:
		r.Resolve(n.Condition)
		r.Resolve(n.Body)
		break

	case *Ast.BinaryExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Right)
		break

	case *Ast.Call:
		r.Resolve(n.Callee)
		for _, arg := range n.Arguments {
			r.Resolve(arg)
		}
		break

	case *Ast.GroupingExpr:
		r.Resolve(n.Expression)
		break

	case *Ast.LiteralExpr:
		r.Resolve(n.Value)
		break

	case *Ast.LogicalExpr:
		r.Resolve(n.Left)
		r.Resolve(n.Right)
		break

	case *Ast.UnaryExpr:
		r.Resolve(n.Right)
		break

	}
}

func (r *Resolver) resolveFunction(stmt *Ast.NamedFunction) {
	enclosingFunction := r.currentFunction
	r.currentFunction = FUNCTION
	r.beginScope()
	for _, arg := range stmt.Params {
		r.declare(arg)
		r.define(arg)
	}
	r.Resolve(stmt.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveLocal(expr Ast.Expr, name *Tokens.Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		scope, _ := r.scopes.Get(i)
		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}

func (r *Resolver) declare(name *Tokens.Token) {
	scope, err := r.scopes.Peek()
	if _, ok := scope[name.Lexeme]; ok {
		Error.ReportParseError(name, "Already a variable with this name in this scope.")
	}
	if err != nil {
		return
	}
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *Tokens.Token) {
	scope, err := r.scopes.Peek()
	if err != nil {
		return
	}
	scope[name.Lexeme] = true
}

func (r *Resolver) beginScope() {
	r.scopes.Push(map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}
