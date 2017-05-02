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
	ctx := map[string]interface{}{
		"id":          1,
		"description": "fool",
	}

	err := Errorc(ctx, msg)
	if got := err.Error(); got != msg {
		t.Errorf(`wrong error message, got "%v", expected "%v"`, got, msg)
		return
	}

	if e := err.(*Error); !reflect.DeepEqual(e.Context, ctx) {
		t.Errorf(`wrong context, got %+v, expected %+v`, e.Context, ctx)
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
	ctx2 := map[string]interface{}{
		"id":          2,
		"description": "bar",
	}
	err2 := Wrapc(err1, ctx2, msg2)

	msg3 := "error message 3"
	ctx3 := map[string]interface{}{
		"id":          3,
		"description": "spam",
	}
	err3 := Wrapc(err2, ctx3, msg3)

	msg4 := "error message 4"
	ctx4 := map[string]interface{}{
		"id":          4,
		"description": "spam",
	}
	err4 := Wrapc(err3, ctx4, msg4)

	got := err4.Error()
	expected := fmt.Sprintf("%s: %s: %s: %s", msg4, msg3, msg2, msg1)
	if got != expected {
		t.Errorf(`wrong error message, got "%s", expected "%s"`, got, expected)
		return
	}
}
