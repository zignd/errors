package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

// Wrapcf returns an error wrapping err, adding contextual information and the provided format specifier.
func Wrapcf(err error, ctx map[string]interface{}, format string, args ...interface{}) error {
	return &Error{
		Message: fmt.Sprintf(format, args...),
		Context: ctx,
		Stack:   callers(),
		Cause:   err,
	}
}

// Stack is an array of program counters that implements fmt.Stringer. Invoke String() in order to obtain a string with the stack trace.
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

// Error is the error struct used internally by the package.
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

func format(err error, lvl int) string {
	t := reflect.TypeOf(err)
	if t != reflect.TypeOf(Error{}) && t != reflect.TypeOf(&Error{}) {
		return fmt.Sprintf("\t%s", err.Error())
	}

	var e Error
	if t.Kind() == reflect.Ptr {
		ep := err.(*Error)
		e = *ep
	} else {
		e = err.(Error)
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("Message:\n\t\"%s\"", e.Message))

	if e.Context != nil {
		b.WriteString("\nContext:")
		for k, v := range e.Context {
			b.WriteString(fmt.Sprintf("\n\t%s: %v", k, v))
		}
	}

	b.WriteString(fmt.Sprintf("\nStack:\n%s", indent(e.Stack.String(), 1)))
	if e.Cause != nil {
		b.WriteString(fmt.Sprintf("\nCause:\n%s", format(e.Cause, 1)))
	}

	return indent(b.String(), lvl)
}

func formatCtx(ctx map[string]interface{}) string {
	return fmt.Sprint(ctx)
}

func indent(s string, times int) string {
	var indent bytes.Buffer
	for i := 0; i < times; i++ {
		indent.WriteString("\t")
	}

	lines := strings.Split(s, "\n")
	nLines := make([]string, len(lines))
	for i, line := range lines {
		nLines[i] = fmt.Sprintf("%s%s", indent.String(), line)
	}

	return strings.Join(nLines, "\n")
}

func callers() *Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st Stack = pcs[0:n]
	return &st
}
