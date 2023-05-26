package Parser

import "fmt"

func (binary *Binary) Print() string {
	return parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (unary *Unary) Print() string {
	return parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (grouping *Grouping) Print() string {
	return parenthesize("group", grouping.Expression)
}

func (literal *Literal) Print() string {
	return fmt.Sprintf("%+v", literal.Value)
}

func (conditional *Conditional) Print() string {
	return fmt.Sprintf("(condition (%s) then (%s) else (%s) )",
		conditional.Condition.Print(),
		conditional.Then.Print(),
		conditional.Else.Print(),
	)
}

func parenthesize(name string, exprs ...Expr) string {
	str := "(" + name
	for _, expr := range exprs {
		str = str + " " + expr.Print()
	}
	str = str + ")"
	return str
}
