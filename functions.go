package errors

import (
	"fmt"
	"runtime"
	"strings"
)

func New(code, message string) error {
	err := &Err{code: code, message: message}
	err.SetLocation(1)
	return err
}

func Errorf(format string, args ...interface{}) error {
	err := &Err{message: fmt.Sprintf(format, args...)}
	err.SetLocation(1)
	return err
}

func getLocation(callDepth int) (string, int) {
	rpc := make([]uintptr, 1)
	n := runtime.Callers(callDepth+2, rpc[:])
	if n < 1 {
		return "", 0
	}
	frame, _ := runtime.CallersFrames(rpc).Next()
	return frame.Function, frame.Line
}

func Trace(other error) error {
	if other == nil {
		return nil
	}

	if err, ok := other.(interface{ Unwrap() []error }); ok {
		curr := &Err{}
		curr.SetLocation(1)
		list := err.Unwrap()
		list = append(list, curr)
		return Join(list...)
	} else if err, ok := other.(*Err); ok {
		curr := &Err{}
		curr.SetLocation(1)
		return Join(err, curr)
	} else {
		prev := New("", other.Error())
		curr := &Err{}
		curr.SetLocation(1)
		return Join(prev, curr)
	}
}

func ErrorStack(err error) string {
	return strings.Join(errorStack(err), "\n")
}

func errorStack(err error) []string {
	if err == nil {
		return nil
	}

	list := make([]error, 0)
	if er, ok := err.(interface{ Unwrap() []error }); ok {
		list = er.Unwrap()
	} else if er, ok := err.(*Err); ok {
		list = append(list, er)
	} else {
		er := New("", err.Error())
		list = append(list, er)
	}

	// We want the first error first
	var lines []string
	for _, e := range list {
		var buff []byte
		if err, ok := e.(Locationer); ok {
			file, line := err.Location()
			// Strip off the leading GOPATH/src path elements.
			if file != "" {
				buff = append(buff, fmt.Sprintf("%s:%d", file, line)...)
				buff = append(buff, ": "...)
			}
		}

		if cerr, ok := e.(*Err); ok {
			code := cerr.Code()
			buff = append(buff, code...)

			message := cerr.Message()
			if len(code) > 0 && len(message) > 0 {
				buff = append(buff, ", "...)
			}

			buff = append(buff, message...)
		} else {
			buff = append(buff, e.Error()...)
		}

		lines = append(lines, string(buff))
		if err == nil {
			break
		}
	}

	// reverse the lines to get the original error, which was at the end of
	// the list, back to the start.
	var result []string
	for i := len(lines); i > 0; i-- {
		result = append(result, lines[i-1])
	}

	return result
}
