package errors

import "testing"

func TestNewMulti(t *testing.T) {
	t.Run("when NewMulti is provided with 3 errors, it should return a new error with the message indicating the number of errors and highlighting the first error", func(t *testing.T) {
		err1 := New("failed 1")
		err2 := New("failed 2")
		err3 := New("failed 3")

		multiErr := NewMulti(err1, err2, err3)
		expected := "first of 3 errors: failed 1"

		if got := multiErr.Error(); got != expected {
			t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
			return
		}
	})
}

func TestAppendMulti(t *testing.T) {
	t.Run("when AppendMulti is provided with a multi error and a new error, it should append the new error to the multi error", func(t *testing.T) {
		multiErr := NewMulti()

		if err := AppendMulti(multiErr, New("failed 1")); err != nil {
			t.Errorf("failed to append an error: %v", err)
			return
		}

		expected := "first of 1 errors: failed 1"

		if got := multiErr.Error(); got != expected {
			t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
			return
		}
	})
}
