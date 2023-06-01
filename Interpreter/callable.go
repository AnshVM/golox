package Interpreter

type LoxCallable struct {
	Arity func() uint
	Call  func(interpreter *Interpreter, arguments []any) any
}
