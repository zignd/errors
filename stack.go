package errors

import (
	"fmt"
	"runtime"
)

// Stack represents a stack trace in the form of a slice of strings.
type Stack []string

// callers returns a stack trace of the calling goroutine.
func callers() Stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(3, pcs[:])
	var st Stack = make([]string, 0)
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])
		if fn == nil {
			st = append(st, "unknown")
		} else {
			file, line := fn.FileLine(pcs[i])
			st = append(st, fmt.Sprintf("%s @ %s:%d", fn.Name(), file, line))
		}
	}
	return st
}
