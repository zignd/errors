package errors

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJSONMarshaling(t *testing.T) {
	t.Run("when marshaling a nested chain of errors.Err errors, should marshal the full chain", func(t *testing.T) {
		err1 := New("context timeout")
		err2 := Wrap(err1, "failed to connect to the database")
		err3 := Wrapd(err2, Data{
			"server": "db-server-01",
		}, "failed to start the server")

		b, err := json.MarshalIndent(err3, "", "  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var errs []map[string]any
		err = json.Unmarshal(b, &errs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(errs) != 3 {
			t.Fatalf("unexpected number of errors, got %d, expected %d", len(errs), 3)
		}

		// testing err3

		if fmt.Sprint(errs[0]["message"]) != err3.(*Err).Message {
			t.Errorf("unexpected error message, got %q, expected %q", errs[0]["message"], err3.(*Err).Message)
		}

		dataErr3, ok := errs[0]["data"]
		if !ok {
			t.Errorf("unexpected data, got undefined key, expected %v", err3.(*Err).Data)
		}

		b1, err := json.Marshal(dataErr3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		b2, err := json.Marshal(err3.(*Err).Data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(b1) != string(b2) {
			t.Errorf("unexpected data, got %s, expected %s", b1, b2)
		}

		// testing err2

		if fmt.Sprint(errs[1]["message"]) != err2.(*Err).Message {
			t.Errorf("unexpected error message, got %q, expected %q", errs[1]["message"], err2.(*Err).Message)
		}

		// testing err1

		if fmt.Sprint(errs[2]["message"]) != err1.(*Err).Message {
			t.Errorf("unexpected error message, got %q, expected %q", errs[2]["message"], err1.(*Err).Message)
		}

		if _, ok := errs[2]["data"]; ok {
			t.Errorf("unexpected data, got %v, expected undefined key", errs[0]["data"])
		}
	})

	t.Run("when marshaling a chain of errors.Err and standard errors, should marshal the full chain", func(t *testing.T) {
		err1 := New("context timeout")
		err2 := fmt.Errorf("failed to connect to the database: %w", err1)
		err3 := Wrap(err2, "failed to start the server")

		b, err := json.MarshalIndent(err3, "", "  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var errs []map[string]any
		err = json.Unmarshal(b, &errs)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(errs) != 3 {
			t.Fatalf("unexpected number of errors, got %d, expected %d", len(errs), 3)
		}

		if fmt.Sprint(errs[0]["message"]) != err3.(*Err).Message {
			t.Errorf("unexpected error message, got %q, expected %q", errs[0]["message"], err3.(*Err).Message)
		}

		if fmt.Sprint(errs[1]["message"]) != err2.Error() {
			t.Errorf("unexpected error message, got %q, expected %q", errs[1]["message"], err2.Error())
		}

		if fmt.Sprint(errs[2]["message"]) != err1.(*Err).Message {
			t.Errorf("unexpected error message, got %q, expected %q", errs[2]["message"], err1.(*Err).Message)
		}
	})
}
