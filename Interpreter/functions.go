package Interpreter

import (
	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Environment"
	"github.com/AnshVM/golox/Error"
)

func CreateFunction(declaration *Ast.Function, closure *Environment.Environment) *LoxCallable {
	Arity := func() uint {
		return uint(len(declaration.Params))
	}
	Call := func(interpreter *Interpreter, arguments []any) (any, error) {
		env := Environment.Environment{Values: map[string]any{}, Enclosing: closure}
		for index, param := range declaration.Params {
			env.Define(param.Lexeme, arguments[index])
		}
		err := interpreter.executeBlock(declaration.Body, &env)
		if err == Error.ErrReturn {
			return interpreter.ReturnValue, nil
		}
		return nil, err
	}
	return &LoxCallable{Arity: Arity, Call: Call}
}
