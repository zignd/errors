package errors

import (
	stderrors "errors"
	"fmt"
	"reflect"
	"testing"
)

// customErr is a custom error type for testing purposes.
type customErr struct {
	msg string
}

// Error returns the error message.
func (c customErr) Error() string { return c.msg }

func TestGo113Compatibility(t *testing.T) {
	t.Run("Wrap should be able to return an error compatible with the standard library Is", func(t *testing.T) {
		// First we create an error using the standard library
		err := stderrors.New("error that gets wrapped")

		// Then we wrap it using our Wrap function
		wrapped := Wrap(err, "wrapped up")

		// Finally we check if the standard library Is function can handle our wrapped error
		if !stderrors.Is(wrapped, err) {
			t.Errorf("Wrap does not support Go 1.13 error chains")
		}
	})

	t.Run("Is should be able to handle errors created and wrapped using the standard Go features", func(t *testing.T) {
		// First we create an error using the standard Go features
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)

		// Then we check if our Is function can handle it
		if !Is(wrapped, err) {
			t.Error("Is failed")
		}
	})

	t.Run("As should be able to handle errors created and wrapped using the standard Go features", func(t *testing.T) {
		// First we create an error using the standard Go features
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)
		target := new(customErr)

		// Then we check if our As function can handle it
		if !As(wrapped, target) {
			t.Error("As failed")
		}
	})

	t.Run("Unwrap should be able to handle errors created and wrapped using the standard Go features", func(t *testing.T) {
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)

		if unwrappedErr := Unwrap(wrapped); !reflect.DeepEqual(unwrappedErr, err) {
			t.Error("Unwrap failed")
		}
	})
}
