package errors

import (
	"encoding/json"
	"fmt"
)

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

func (e *Err) MarshalJSON() ([]byte, error) {
	return json.Marshal(toMapsSlice(e))
}
