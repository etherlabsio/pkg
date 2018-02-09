package errors

import (
	"encoding/json"
	"testing"
)

func Test_MarshalJSON(t *testing.T) {
	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := New(Op("Get"), IO, "network unreachable")
	e1 = Wrap(e1, "something went wrong")
	e1 = Wrap(e1, "i don't know what")

	// Nested error.
	e2 := New(Op("Read"), Other, e1)
	e2 = Wrap(e2, "nothing is wrong")

	b, _ := json.Marshal(e2)

	in := e2.(*Error)
	out := new(Error)
	json.Unmarshal(b, out)

	t.Logf("err2 err: %s, err3 err: %s", in.Error(), out.Error())
	// Compare elementwise.
	if in.Op != out.Op {
		t.Errorf("expected Op %q; got %q", in.Op, out.Op)
	}
	if in.Kind != out.Kind {
		t.Errorf("expected kind %d; got %d", in.Kind, out.Kind)
	}
	// Note that error will have lost type information, so just check its Error string.
	if in.Error() != out.Error() {
		t.Errorf("expected Err => %s; got => %s", in, out)
	}
}
