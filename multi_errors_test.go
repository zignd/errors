package errors

import "testing"

func TestNewMulti(t *testing.T) {
	err1 := New("failed 1")
	err2 := New("failed 2")
	err3 := New("failed 3")

	multiErr := NewMulti(err1, err2, err3)
	expected := "first of 3 errors: failed 1"

	if got := multiErr.Error(); got != expected {
		t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
		return
	}
}

func TestAppendMulti(t *testing.T) {
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
}
