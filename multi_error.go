package errors

import "fmt"

// NewMulti returns a new errors.MultiError with the provided errs.
func NewMulti(errs ...error) error {
	return &MultiError{
		Errors: errs,
	}
}

// AppendMulti appends err to multi.Errors and returns an error if multi is not of the type errors.MultiError.
func AppendMulti(multi error, err error) error {
	m, ok := multi.(*MultiError)
	if !ok {
		return New("the multi parameter is supposed to be given an errors.MultiError value")
	}

	m.Errors = append(m.Errors, err)

	return nil
}

// MultiError is the error struct for multiple errors used internally by the package. This type should be only be used for type assertions.
type MultiError struct {
	Errors []error
}

func (m MultiError) Error() string {
	if len(m.Errors) == 0 {
		return ""
	}

	return fmt.Sprintf("first of %d errors: %s", len(m.Errors), m.Errors[0].Error())
}
