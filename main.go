package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Interpreter"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Scanner"
)

func run(source string) {
	scanner := Scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := Parser.NewParser(tokens)
	expr := parser.Parse()
	if Error.HadError || Error.HadRuntimeError {
		return
	}
	Interpreter.Interpret(expr)
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
	if Error.HadRuntimeError {
		os.Exit(70)
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
		Error.HadRuntimeError = false
	}
}

func main() {
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
