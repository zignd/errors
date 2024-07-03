package errors

import "fmt"

// Err is the error struct used internally by the package. This type should only be used for type assertions.
type Err struct {
	Message string `json:"message"`
	Data    Data   `json:"data,omitempty"`
	Stack   Stack  `json:"stack"`
	Cause   error  `json:"cause,omitempty"`
}

func (e Err) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}

	return e.Message
}

// Format implements fmt.Formatter. It only accepts the '+v' and 's' formats.
func (e Err) Format(s fmt.State, verb rune) {
	if verb == 'v' && s.Flag('+') {
		fmt.Fprintf(s, "%s", format(e, 0))
	} else {
		fmt.Fprintf(s, "%s", e.Error())
	}
}

func (e Err) Unwrap() error {
	return e.Cause
}

// WithStack adds a stack trace to the provided error if it is an Err or *Err.
func WithStack(err error) error {
	if e, ok := err.(Err); ok {
		e.Stack = callers()
		return e
	} else if e, ok := err.(*Err); ok {
		e.Stack = callers()
		return e
	} else {
		return err
	}
}
