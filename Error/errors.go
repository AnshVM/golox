package Error

import "fmt"

var HadError = false

func Report(line uint, where string, message string) {
	fmt.Println("[line " + fmt.Sprint(line) + "] Error" + where + ": " + message)
	HadError = true
}

func Error(line uint, message string) {
	Report(line, "", message)
}
