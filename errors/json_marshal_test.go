package errors

import (
	"bytes"
	"encoding/json"
	"testing"
)

func testJSONRoundTrip(err error) error {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(err)
	var e Error
	json.NewDecoder(b).Decode(&e)
	return &e
}

func TestErrorMarshalJSON(t *testing.T) {
	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := New("network unreachable", IO, Str("tcp connection dropped"))
	input := e1.(*Error)
	output, ok := testJSONRoundTrip(e1).(*Error)
	if !ok {
		t.Errorf("expected pointer err, but failed")
	}

	if output.kind != input.kind {
		t.Errorf("expected kind %d; got %d", input.kind, output.kind)
	}

	// Note that error will have lost type information, so just check its Error string.
	if input.Error() != output.Error() {
		t.Errorf("expected Err => %s; got => %s", input, output)
	}
}

func TestWrappedMarshalJSON(t *testing.T) {
	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := New("network unreachable", IO, Str("tcp connection dropped"))
	op := Op("TestWrappedMarshalJSON")
	e2 := WithOp(e1, op)

	err := Unwrap(e2)

	input := err.(*Error)
	output, ok := testJSONRoundTrip(err).(*Error)
	if !ok {
		t.Errorf("expected pointer err, but failed")
	}

	if output.kind != input.kind {
		t.Errorf("expected kind %d; got %d", input.kind, output.kind)
	}

	// Note that error will have lost type information, so just check its Error string.
	if input.Error() != output.Error() {
		t.Errorf("expected Err => %s; got => %s", input, output)
	}
}
