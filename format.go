package errors

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

// format returns a formatted string representation of the error and its cause.
func format(err error, lvl int) string {
	t := reflect.TypeOf(err)
	if t != reflect.TypeOf(Err{}) && t != reflect.TypeOf(&Err{}) {
		return fmt.Sprintf("\t%s", err.Error())
	}

	var e Err
	if t.Kind() == reflect.Ptr {
		ep := err.(*Err)
		e = *ep
	} else {
		e = err.(Err)
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("message:\n\t\"%s\"", e.Message))

	if e.Data != nil {
		b.WriteString("\ndata:")
		for k, v := range e.Data {
			b.WriteString(fmt.Sprintf("\n\t%s: %v", k, v))
		}
	}

	if e.Stack != nil && len(e.Stack) > 0 {
		firstStackLine := e.Stack[0]
		b.WriteString(fmt.Sprintf("\nstack:\n%s", indent(firstStackLine, 1)))
		for i := 1; i < len(e.Stack); i++ {
			b.WriteString(fmt.Sprintf("\n%s", indent(e.Stack[i], 1)))
		}
	}

	if e.Cause != nil {
		b.WriteString(fmt.Sprintf("\ncause:\n%s", format(e.Cause, 1)))
	}

	return indent(b.String(), lvl)
}

// indent indents a string by the given number of times.
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
