package Error

import "errors"

var (
	CANNOT_READ_FILE = "CANNOT_READ_FILE"
	ErrParseError    = errors.New("ErrParseError")
)
