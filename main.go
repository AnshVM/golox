package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/AnshVM/golox/Error"
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

func main() {
	if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}
