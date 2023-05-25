package Error

import (
	"fmt"

	"github.com/AnshVM/golox/Tokens"
)

var HadError = false

func Report(line uint, where string, message string) {
	fmt.Println("[line " + fmt.Sprint(line) + "] Error" + where + ": " + message)
	HadError = true
}

func ReportParseError(token *Tokens.Token, message string) {
	if token.Type == Tokens.EOF {
		Report(token.Line, "at end", message)
	} else {
		Report(token.Line, fmt.Sprintf("at '%s'", token.Lexeme), message)
	}
}

func ReportScanError(line uint, message string) {
	Report(line, "", message)
}
