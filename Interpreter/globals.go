package Interpreter

import "time"

func Clock() *LoxCallable {
	Call := func(_ *Interpreter, _ []any) (any, error) {
		return time.Now().Second(), nil
	}
	Arity := func() uint {
		return 0
	}
	return &LoxCallable{Call: Call, Arity: Arity}
}
