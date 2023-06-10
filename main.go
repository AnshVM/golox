package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AnshVM/golox/Environment"
	"github.com/AnshVM/golox/Error"
	"github.com/AnshVM/golox/Interpreter"
	"github.com/AnshVM/golox/Parser"
	"github.com/AnshVM/golox/Resolver"
	"github.com/AnshVM/golox/Scanner"
)

func run(i *Interpreter.Interpreter, source string) {
	scanner := Scanner.NewScanner(source)
	tokens := scanner.ScanTokens()

	parser := Parser.NewParser(tokens)
	stmts := parser.Parse()
	resolver := Resolver.NewResolver(i)
	resolver.Resolve(stmts)
	if Error.HadError {
		return
	}
	i.Interpret(stmts)
}

func runFile(i *Interpreter.Interpreter, path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.New(Error.CANNOT_READ_FILE)
	}
	run(i, string(data))
	if Error.HadError {
		os.Exit(65)
	}
	if Error.HadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func runPrompt(i *Interpreter.Interpreter) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		line = line[:len(line)-1]
		run(i, line)
		Error.HadError = false
		Error.HadRuntimeError = false
	}
}

func main() {
	globals := Environment.Environment{Values: make(map[string]any)}
	interpreter := Interpreter.NewInterpreter(&globals)
	if len(os.Args) == 2 {
		runFile(interpreter, os.Args[1])
	} else {
		runPrompt(interpreter)
	}
}
