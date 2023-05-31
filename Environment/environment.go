package Environment

import (
	"fmt"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

type Environment struct {
	Values    map[string]any
	Enclosing *Environment
}

func (env *Environment) Define(name *Tokens.Token, value any) {
	env.Values[name.Lexeme] = value
}

func (env *Environment) Get(name *Tokens.Token) (any, error) {
	if val, ok := env.Values[name.Lexeme]; ok {
		return val, nil
	} else if env.Enclosing != nil {
		return env.Enclosing.Get(name)
	}
	Error.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
	return nil, Error.ErrRuntimeError

}

func (env *Environment) Assign(name *Tokens.Token, value any) error {
	if _, ok := env.Values[name.Lexeme]; ok {
		env.Values[name.Lexeme] = value
	} else if env.Enclosing != nil {
		env.Enclosing.Assign(name, value)
	} else {
		Error.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
		return Error.ErrRuntimeError
	}
	return nil
}
