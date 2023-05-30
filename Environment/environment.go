package Environment

import (
	"fmt"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Tokens"
)

type Environment struct {
	Values map[string]any
}

func (env *Environment) Define(name *Tokens.Token, value any) {
	env.Values[name.Lexeme] = value
}

func (env *Environment) Get(name *Tokens.Token) any {
	if val, ok := env.Values[name.Lexeme]; ok {
		return val
	}
	Error.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
	return nil
}

func (env *Environment) Assign(name *Tokens.Token, value any) {
	if _, ok := env.Values[name.Lexeme]; ok {
		env.Values[name.Lexeme] = value
	} else {
		Error.ReportRuntimeError(name, fmt.Sprintf("Undefined variable %s", name.Lexeme))
	}
}
