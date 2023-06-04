package Interpreter

import (
	"github.com/AnshVM/golox/Ast"
	"github.com/AnshVM/golox/Environment"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

func CreateFunctionCallable(body []Ast.Stmt, params []*Tokens.Token, closure *Environment.Environment) *LoxCallable {
	Arity := func() uint {
		return uint(len(params))
	}
	Call := func(interpreter *Interpreter, arguments []any) (any, error) {
		env := Environment.Environment{Values: map[string]any{}, Enclosing: closure}
		for index, param := range params {
			env.Define(param.Lexeme, arguments[index])
		}
		err := interpreter.executeBlock(body, &env)
		if err == Error.ErrReturn {
			return interpreter.ReturnValue, nil
		}
		return nil, err
	}
	return &LoxCallable{Arity: Arity, Call: Call}
}
