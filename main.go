package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Scanner"
)

func run(source string) {
	scanner := Scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
}

func runFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(Error.CANNOT_READ_FILE)
	}
	run(string(data))
	if Error.HadError {
		os.Exit(65)
	}
	return nil
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1]
		run(line)
		Error.HadError = false
	}
}

func PrintExpr() {
	expr := Parser.Binary{
		Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1),
		Left:     &Parser.Unary{Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1), Right: &Parser.Literal{Value: 123}},
		Right:    &Parser.Grouping{Expression: &Parser.Literal{Value: 45.67}},
	}

	expr2 := Parser.Binary{
		Operator: Scanner.NewToken(Scanner.STAR, "*", nil, 1),
		Left: &Parser.Grouping{
			Expression: &Parser.Binary{
				Operator: Scanner.NewToken(Scanner.PLUS, "+", nil, 1),
				Right:    &Parser.Literal{Value: 1},
				Left:     &Parser.Literal{Value: 2},
			},
		},
		Right: &Parser.Grouping{
			Expression: &Parser.Binary{
				Operator: Scanner.NewToken(Scanner.MINUS, "-", nil, 1),
				Right:    &Parser.Literal{Value: 4},
				Left:     &Parser.Literal{Value: 3},
			},
		},
	}
	fmt.Println(expr.Print())
	fmt.Println(expr.PrintRPN())
	fmt.Println(expr2.PrintRPN())
	fmt.Println(expr2.Print())
}

func main() {
	// if len(os.Args) == 2 {
	// 	runFile(os.Args[1])
	// } else {
	// 	runPrompt()
	// }
	PrintExpr()
}
