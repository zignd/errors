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
	t.Run("when Wrap is used to wrap a standard error, it should return an error compatible with the standard library Is", func(t *testing.T) {
		// First we create an error using the standard library
		err := stderrors.New("error that gets wrapped")

		// Then we wrap it using our Wrap function
		wrapped := Wrap(err, "wrapped up")

		// Finally we check if the standard library Is function can handle our wrapped error
		if !stderrors.Is(wrapped, err) {
			t.Errorf("our Wrap does not support Go 1.13 error chains")
		}
	})

	// Is should be able to handle errors created and wrapped using the standard Go features
	t.Run("when Is is used to check if an error is a certain error, it should behave just like the equivalent Is function in the standard library", func(t *testing.T) {
		// First we create an error using the standard Go features
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)

		// Then we check if our Is function can handle it
		if !Is(wrapped, err) {
			t.Error("Is failed")
		}

		// Finally just to make sure, we check if the standard library Is function can handle it
		if !stderrors.Is(wrapped, err) {
			t.Error("stderrors.Is failed")
		}
	})

	t.Run("when As is used to check if an error is a certain error, it should behave just like the equivalent As function in the standard library", func(t *testing.T) {
		// First we create an error using the standard Go features
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)
		target := new(customErr)

		// Then we check if our As function can handle it
		if !As(wrapped, target) {
			t.Error("As failed")
		}

		// Finally just to make sure, we check if the standard library As function can handle it
		if !stderrors.As(wrapped, target) {
			t.Error("stderrors.As failed")
		}
	})

	// Unwrap should be able to handle errors created and wrapped using the standard Go features
	t.Run("when Unwrap is used to unwrap an error, it should behave just like the equivalent Unwrap function in the standard library", func(t *testing.T) {
		err := customErr{msg: "test message"}
		wrapped := fmt.Errorf("wrap it: %w", err)

		if unwrappedErr := Unwrap(wrapped); !reflect.DeepEqual(unwrappedErr, err) {
			t.Error("Unwrap failed")
		}

		if unwrappedErr := stderrors.Unwrap(wrapped); !reflect.DeepEqual(unwrappedErr, err) {
			t.Error("stderrors.Unwrap failed")
		}
	})
}
