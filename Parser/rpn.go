package Parser

import "fmt"

func (literal *Literal) PrintRPN() string {
	return fmt.Sprintf("%v", literal.Value)
}

func (unary *Unary) PrintRPN() string {
	return fmt.Sprintf("%s%s", unary.Operator.Lexeme, unary.Right.PrintRPN())
}

func (binary *Binary) PrintRPN() string {
	return fmt.Sprintf("%s %s %s",
		binary.Left.PrintRPN(),
		binary.Right.PrintRPN(),
		binary.Operator.Lexeme,
	)
}

func (grouping *Grouping) PrintRPN() string {
	return grouping.Expression.PrintRPN()
}

//  -1 + 2 -> 1- 2 +
