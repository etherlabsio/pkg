package errors

import "testing"

func TestUnWrap(t *testing.T) {
	const op = Op("TestKindOf")
	e0 := Str("tcp connection dropped")
	e1 := New("network unreachable", IO, e0)
	e2 := WithMessage(e1, "annotation0")
	e3 := WithOpf(e2, op, "annotation1")

	want := e1
	if got := Unwrap(e3); got != want {
		t.Errorf("KindOf() = %v, want %v", got, want)
	}

	// add second level kind
	e4 := WithKind(e3, Permission, "permission error")
	e5 := WithMessage(e4, "annotation")
	want = e4
	if got := Unwrap(e5); got != want {
		t.Errorf("KindOf() = %v, want %v", got, want)
	}
}

func TestKindOf(t *testing.T) {
	const op = Op("TestKindOf")
	e0 := Str("tcp connection dropped")
	e1 := New("network unreachable", IO, e0)
	e2 := WithMessage(e1, "annotation0")
	e3 := WithOpf(e2, op, "annotation1")

	want := IO
	if got := KindOf(e3); got != want {
		t.Errorf("KindOf() = %v, want %v", got, want)
	}

	// add second level kind
	e4 := WithKind(e3, Permission, "permission error")
	e5 := WithMessage(e4, "annotation")
	want = Permission
	if got := KindOf(e5); got != want {
		t.Errorf("KindOf() = %v, want %v", got, want)
	}

	want = Internal
	if got := KindOf(e0); got != want {
		t.Errorf("when Error was not present: KindOf() = %v, want %v", got, want)
	}
}
