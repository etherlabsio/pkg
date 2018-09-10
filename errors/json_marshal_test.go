package errors

import (
	"testing"
)

func Test_MarshalJSON(t *testing.T) {

	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := New("network unreachable", Op("Get"), IO)
	t.Logf("err2 err: %s, err3 err: %s", in.Error(), out.Error())

	if in.Kind != out.Kind {
		t.Errorf("expected kind %d; got %d", in.Kind, out.Kind)
	}
	// Note that error will have lost type information, so just check its Error string.
	if in.Error() != out.Error() {
		t.Errorf("expected Err => %s; got => %s", in, out)
	}
}
