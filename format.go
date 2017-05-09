package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

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
