package errors

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	msg := "error message"
	if got := New(msg).Error(); got != msg {
		t.Errorf(`wrong error message, got "%v", expected "%v"`, got, msg)
		return
	}
}

func TestErrorc(t *testing.T) {
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

	if e := err.(*Error); !reflect.DeepEqual(e.Data, data) {
		t.Errorf(`wrong data, got %+v, expected %+v`, e.Data, data)
		return
	}
}

func TestWrap(t *testing.T) {
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
}

func TestWrapc(t *testing.T) {
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
}
