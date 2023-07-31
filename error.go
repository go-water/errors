package errors

import (
	stderrors "errors"
	"fmt"
	"strings"
)

type Err struct {
	code     string
	message  string
	function string
	line     int
}

type Locationer interface {
	Location() (function string, line int)
}

func (e *Err) Location() (function string, line int) {
	return e.function, e.line
}

func (e *Err) Message() string {
	return e.message
}

func (e *Err) Code() string {
	return e.code
}

func (e *Err) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *Err) SetLocation(callDepth int) {
	e.function, e.line = getLocation(callDepth + 1)
}

func (e *Err) StackTrace() []string {
	return errorStack(e)
}

func (e *Err) ErrorType() string {
	if e == nil || len(e.code) == 0 {
		return ""
	}

	return strings.Split(e.code, ".")[0]
}

func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
