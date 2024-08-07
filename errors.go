package errors

import (
	"fmt"
	"reflect"
)

type Data map[string]any

// New returns an error with the provided message.
func New(msg string) error {
	return &Err{
		Message: msg,
		Stack:   callers(),
	}
}

// Errord returns an error with additional data and the provided message.
func Errord(data Data, msg string) error {
	return &Err{
		Message: msg,
		Data:    data,
		Stack:   callers(),
	}
}

// Errorf returns an error with the provided format specifier.
func Errorf(format string, args ...any) error {
	return &Err{
		Message: fmt.Sprintf(format, args...),
		Stack:   callers(),
	}
}

// Errordf returns an error with additional data and the provided format specifier.
func Errordf(data Data, format string, args ...any) error {
	return &Err{
		Message: fmt.Sprintf(format, args...),
		Data:    data,
		Stack:   callers(),
	}
}

// Wrap returns an error wrapping err and adding the provided format specifier.
func Wrap(err error, msg string) error {
	return &Err{
		Message: msg,
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapd returns an error wrapping err, adding additional data and the provided message.
func Wrapd(err error, data Data, msg string) error {
	return &Err{
		Message: msg,
		Data:    data,
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapf returns an error wrapping err and adding the provided format specifier.
func Wrapf(err error, format string, args ...any) error {
	return &Err{
		Message: fmt.Sprintf(format, args...),
		Stack:   callers(),
		Cause:   err,
	}
}

// Wrapdf returns an error wrapping err, adding additional data and the provided format specifier.
func Wrapdf(err error, data Data, format string, args ...any) error {
	return &Err{
		Message: fmt.Sprintf(format, args...),
		Data:    data,
		Stack:   callers(),
		Cause:   err,
	}
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

// WithCause adds a cause to the provided error if it is an Err or *Err.
func WithCause(err error, cause error) error {
	if e, ok := err.(Err); ok {
		e.Cause = cause
		return e
	} else if e, ok := err.(*Err); ok {
		e.Cause = cause
		return e
	} else {
		return err
	}
}

// IsErrComposition returns true if the provided error is a composition of Err or *Err.
func IsErrComposition(err error) bool {
	typeOfErr := reflect.TypeOf(err)

	if typeOfErr.Kind() == reflect.Pointer {
		typeOfErr = typeOfErr.Elem()
	}

	if typeOfErr.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < typeOfErr.NumField(); i++ {
		if typeOfErr.Field(i).Type == reflect.TypeOf(Err{}) {
			return true
		}

		if typeOfErr.Field(i).Type == reflect.TypeOf((*Err)(nil)) {
			return true
		}
	}

	return false
}
