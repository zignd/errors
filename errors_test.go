package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// CustomError is a custom error type composed with Err.
type CustomError struct {
	*Err
}

// NewCustomError returns a new CustomError and adds a stack trace.
func NewCustomError(message string) error {
	customError := CustomError{Err: &Err{Message: message}}
	WithStack(customError.Err)
	return customError
}

// CustomError2 is a custom error type composed with Err.
type CustomError2 struct {
	*Err
}

// NewCustom2Error returns a new CustomError2 and adds a cause to the error.
func NewCustom2Error(message string, cause error) error {
	customError2 := CustomError2{Err: &Err{Message: message}}
	WithCause(customError2.Err, cause)
	return customError2
}

// customError3 is a custom error type not composed with Err.
type customError3 struct{}

func (e customError3) Error() string {
	return "this is a custom error type"
}

func TestNew(t *testing.T) {
	t.Run("when New is provided with a message, it should create a new error with the message", func(t *testing.T) {
		msg := "error message"
		if got := New(msg).Error(); got != msg {
			t.Errorf(`wrong error message, got "%v", expected "%v"`, got, msg)
			return
		}
	})
}

func TestErrord(t *testing.T) {
	t.Run("when Errord is provided with a message and data, it should create a new error with the message and data", func(t *testing.T) {
		msg := "error message"
		data := Data{
			"id":          1,
			"description": "fool",
		}

		err := Errord(data, msg)
		if got := err.Error(); got != msg {
			t.Errorf(`wrong error message, got "%v", expected "%v"`, got, msg)
			return
		}

		if e := err.(*Err); !reflect.DeepEqual(e.Data, data) {
			t.Errorf(`wrong data, got "%+v", expected "%+v"`, e.Data, data)
			return
		}
	})
}

func TestWrap(t *testing.T) {
	t.Run("when Wrap is provided with an error and a message, it should create a new error with the message and the provided error as the cause", func(t *testing.T) {
		msg1 := "error message 1"
		err1 := New(msg1)
		msg2 := "error message 2"
		err2 := Wrap(err1, msg2)
		msg3 := "error message 3"
		err3 := Wrap(err2, msg3)
		got := err3.Error()
		expected := fmt.Sprintf("%s: %s: %s", msg3, msg2, msg1)
		if got != expected {
			t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
			return
		}
	})
}

func TestWrapc(t *testing.T) {
	t.Run("when Wrapd is provided with an error and data, it should add the data to the error", func(t *testing.T) {
		msg1 := "error message 1"
		err1 := errors.New(msg1)

		msg2 := "error message 2"
		data2 := Data{
			"id":          2,
			"description": "bar",
		}
		err2 := Wrapd(err1, data2, msg2)

		msg3 := "error message 3"
		data3 := Data{
			"id":          3,
			"description": "spam",
		}
		err3 := Wrapd(err2, data3, msg3)

		msg4 := "error message 4"
		data4 := Data{
			"id":          4,
			"description": "spam",
		}
		err4 := Wrapd(err3, data4, msg4)

		got := err4.Error()
		expected := fmt.Sprintf("%s: %s: %s: %s", msg4, msg3, msg2, msg1)
		if got != expected {
			t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
			return
		}

		if e := err4.(*Err); !reflect.DeepEqual(e.Data, data4) {
			t.Errorf(`wrong data, got "%+v", expected "%+v"`, e.Data, data4)
			return
		}

		if e := err3.(*Err); !reflect.DeepEqual(e.Data, data3) {
			t.Errorf(`wrong data, got "%+v", expected "%+v"`, e.Data, data3)
			return
		}

		if e := err2.(*Err); !reflect.DeepEqual(e.Data, data2) {
			t.Errorf(`wrong data, got "%+v", expected "%+v"`, e.Data, data2)
			return
		}
	})
}

func TestWithStack(t *testing.T) {
	t.Run("when WithStack is provided with an error of type Err, it should add a stack trace to the error", func(t *testing.T) {
		err := NewCustomError("this is a custom error type with stack")

		if err.(CustomError).Stack == nil {
			t.Fatal("expected stack to be not nil, got nil")
			return
		}

		outputStr := fmt.Sprintf("%+v", err)
		if !strings.Contains(outputStr, "message:") {
			t.Errorf(`expected "message:" to be in the output string, got "%v"`, outputStr)
		}
		if !strings.Contains(outputStr, "stack:") {
			t.Errorf(`expected "stack:" to be in the output string, got "%v"`, outputStr)
		}
	})
}

func TestWithCause(t *testing.T) {
	t.Run("when WithCause is provided with an error and a cause, it should add the cause to the error", func(t *testing.T) {
		causeErr := New("inner error")
		err := NewCustom2Error("outer error", causeErr)

		if err.(CustomError2).Cause != causeErr {
			t.Errorf(`expected cause to be "%v", got "%v"`, causeErr, err.(CustomError2).Cause)
		}
	})
}

func TestIsErrComposition(t *testing.T) {
	t.Run("when a custom error type is composed with *Err, it should return true", func(t *testing.T) {
		err := NewCustomError("this is a custom error type with stack")
		if !IsErrComposition(err) {
			t.Errorf("expected IsErrComposition to return true, got false")
		}
	})

	t.Run("when an error type is Pointer but the element type is a struct not composed with *Err, it should return false", func(t *testing.T) {
		err := errors.New("this is a regular error")
		if IsErrComposition(err) {
			t.Errorf("expected IsErrComposition to return false, got true")
		}
	})

	t.Run("when a custom error type is not composed with *Err, it should return false", func(t *testing.T) {
		err := customError3{}
		if IsErrComposition(err) {
			t.Errorf("expected IsErrComposition to return false, got true")
		}
	})
}
