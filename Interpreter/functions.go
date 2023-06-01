package Interpreter

import (
	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Environment"
)

func CreateFunction(declaration *Ast.Function) *LoxCallable {
	Arity := func() uint {
		return uint(len(declaration.Params))
	}
	Call := func(interpreter *Interpreter, arguments []any) (any, error) {
		env := Environment.Environment{Values: map[string]any{}, Enclosing: interpreter.Env}
		for index, param := range declaration.Params {
			env.Define(param.Lexeme, arguments[index])
		}
		return interpreter.executeBlock(declaration.Body, &env)
	}
	return &LoxCallable{Arity: Arity, Call: Call}
}
