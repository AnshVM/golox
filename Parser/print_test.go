package Parser

import (
	"testing"

	"github.com/AnshVM/golox/Scanner"
	"github.com/stretchr/testify/assert"
)

var expr1 = Binary{
	Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1),
	Left:     &Unary{Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1), Right: &Literal{Value: 123}},
	Right:    &Grouping{Expression: &Literal{Value: 45.67}},
}

var expr2 = Binary{
	Operator: Scanner.NewToken(Scanner.STAR, "*", nil, 1),
	Left: &Grouping{
		Expression: &Binary{
			Operator: Scanner.NewToken(Scanner.PLUS, "+", nil, 1),
			Right:    &Literal{Value: 1},
			Left:     &Literal{Value: 2},
		},
	},
	Right: &Grouping{
		Expression: &Binary{
			Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1),
			Right:    &Literal{Value: 4},
			Left:     &Literal{Value: 3},
		},
	},
}

func TestPrint(t *testing.T) {
	assert.Equal(t, "(- (- 123) (group 45.67))", expr1.Print())
	assert.Equal(t, "(* (group (+ 2 1)) (group (- 3 4)))", expr2.Print())
}

func TestPrintRPN(t *testing.T) {
	assert.Equal(t, "-123 45.67 -", expr1.PrintRPN())
	assert.Equal(t, "2 1 + 3 4 - *", expr2.PrintRPN())
}
