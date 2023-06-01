package Interpreter

import "time"

func Clock() *LoxCallable {
	Call := func(_ *Interpreter, _ []any) any {
		return time.Now().Second()
	}
	Arity := func() uint {
		return 0
	}
	return &LoxCallable{Call: Call, Arity: Arity}
}
