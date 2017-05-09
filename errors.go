package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

// New returns an error with the provided message.
func New(msg string) error {
	return &Error{
		Message: msg,
		Stack:   callers(),
	}
}

// Errorc returns an error with contextual information and the provided message.
func Errorc(ctx map[string]interface{}, msg string) error {
	return &Error{
		Message: msg,
		Context: ctx,
		Stack:   callers(),
	}
}

// Errorf returns an error with the provided format specifier.
func Errorf(format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Stack:   callers(),
	}
}

// Errorcf returns an error with contextual information and the provided format specifier.
func Errorcf(ctx map[string]interface{}, format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Context: ctx,
		Stack:   callers(),
	}
}

// Wrap returns an error wrapping err and adding the provided format specifier.
func Wrap(err error, msg string) error {
	return &Error{
		Message: msg,
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapc returns an error wrapping err, adding contextual information and the provided message.
func Wrapc(err error, ctx map[string]interface{}, msg string) error {
	return &Error{
		Message: msg,
		Context: ctx,
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapf returns an error wrapping err and adding the provided format specifier.
func Wrapf(err error, format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapcf returns an error wrapping err, adding contextual information and the provided format specifier.
func Wrapcf(err error, ctx map[string]interface{}, format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Context: ctx,
		Stack:   callers(),
		Cause:   err,
	}
}

// Stack is an array of program counters that implements fmt.Stringer. Call String() in order to obtain a string with the stack trace.
type Stack []uintptr

func (s *Stack) String() string {
	var b bytes.Buffer
	for i := 0; i < len(*s); i++ {
		if i != 0 {
			b.WriteString("\n")
		}
		pc := (*s)[i]
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			b.WriteString("unknown")
		} else {
			file, line := fn.FileLine(pc)
			b.WriteString(fmt.Sprintf("%s\n\t%s:%d", fn.Name(), file, line))
		}
	}
	return b.String()
}

// Error is the error struct used internally by the package. This type should only be used for type assertions.
type Error struct {
	Message string
	Context map[string]interface{}
	Stack   *Stack
	Cause   error
}

func (e Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}

	return e.Message
}

// Format implements fmt.Formatter. It only accepts the '+v' and 's' formats.
func (e Error) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprintf(s, "%s", format(e, 0))
	} else {
		fmt.Fprintf(s, "%s", e.Error())
	}
}
