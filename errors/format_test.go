package errors

import "testing"

func TestWrap(t *testing.T) {
	err := New(IO, "network issue", Op("get"))
	err = Wrap(err, "figure it out")
	err = Wrap(err, "don't worry about it")

	want := "get: 3: don't worry about it: figure it out: network issue"
	if err.Error() != want {
		t.Errorf("Wrap() error = %s, wantErr %s", err, want)
	}
	e := err.(*Error)
	if e.Kind != IO {
		t.Errorf("Wrap() expected kind = %d, got = %d", IO, e.Kind)
	}
	if e.Op != "get" {
		t.Errorf("Wrap() expected op = %s, got = %s", "get", e.Op)
	}

	err = nil
	nilErr := Wrap(err, "something")
	if nilErr != nil {
		t.Errorf("Wrap() nilErr = %v, wanted nil", err)
	}
}

func TestWrapf(t *testing.T) {
	err := New(IO, "network issue", Op("get"))
	err = Wrapf(err, "figure it out %d", 1)
	err = Wrapf(err, "don't worry about it %s", "buddy")

	want := "get: 3: don't worry about it buddy: figure it out 1: network issue"
	if err.Error() != want {
		t.Errorf("Wrapf() error = %s, wantErr %s", err, want)
	}
	e := err.(*Error)
	if e.Kind != IO {
		t.Errorf("Wrapf() expected kind = %d, got = %d", IO, e.Kind)
	}
	if e.Op != "get" {
		t.Errorf("Wrapf() expected op = %s, got = %s", "get", e.Op)
	}

	err = nil
	nilErr := Wrap(err, "something")
	if nilErr != nil {
		t.Errorf("Wrapf() nilErr = %v, wanted nil", err)
	}
}

func TestWithMessage(t *testing.T) {
	err := New(IO, "network issue", Op("get"))
	err = WithMessage(err, "figure it out")
	err = WithMessage(err, "don't worry about it")

	want := "get: 3: don't worry about it: figure it out: network issue"
	if err.Error() != want {
		t.Errorf("WithMessage() error = %s, wantErr %s", err, want)
	}
	e := err.(*Error)
	if e.Kind != IO {
		t.Errorf("WithMessage() expected kind = %d, got = %d", IO, e.Kind)
	}
	if e.Op != "get" {
		t.Errorf("WithMessage() expected op = %s, got = %s", "get", e.Op)
	}

	err = nil
	nilErr := WithMessage(err, "something")
	if nilErr != nil {
		t.Errorf("WithMessage() nilErr = %v, wanted nil", err)
	}
}
