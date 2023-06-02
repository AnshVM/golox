package Error

import "errors"

var (
	CANNOT_READ_FILE = "CANNOT_READ_FILE"
	ErrParseError    = errors.New("ParseError")
	ErrRuntimeError  = errors.New("RuntimeError")
	ErrReturn        = errors.New("Return")
)
