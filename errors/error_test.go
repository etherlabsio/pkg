package errors

import (
	"io"
	"testing"
)

func TestErrorf(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{Errorf("read error without format specifiers"), "read error without format specifiers"},
		{Errorf("read error with %d format specifier", 1), "read error with 1 format specifier"},
	}

	for _, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("Errorf(%v): got: %q, want %q", tt.err, got, tt.want)
		}
	}
}

func TestDoesNotChangePreviousError(t *testing.T) {
	err := New(Permission)
	err2 := New(Op("I will NOT modify err"), err)

	expected := "I will NOT modify err: 2"
	if err2.Error() != expected {
		t.Fatalf("Expected %q, got %q", expected, err2)
	}
	kind := err.(*Error).Kind
	if kind != Permission {
		t.Fatalf("Expected kind %v, got %v", Permission, kind)
	}
}

func TestSeparator(t *testing.T) {
	defer func(prev string) {
		Separator = prev
	}(Separator)
	Separator = ":: "

	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := New(Op("Get"), IO, "network unreachable")

	// Nested error.
	e2 := New(Op("Read"), Other, e1)

	want := "Read: 3:: Get: network unreachable"
	if got := e2.Error(); got != want {
		t.Errorf("expected %q; got %q", want, got)
	}
}

func TestNoArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = New()
}

type matchTest struct {
	err1, err2 error
	matched    bool
}

const (
	op  = Op("Op")
	op1 = Op("Op1")
	op2 = Op("Op2")
)

var matchTests = []matchTest{
	// Errors not of type *Error fail outright.
	{nil, nil, false},
	{io.EOF, io.EOF, false},
	{New(io.EOF), io.EOF, false},
	{io.EOF, New(io.EOF), false},
	// Success. We can drop fields from the first argument and still match.
	{New(io.EOF), New(io.EOF), true},
	{New(op, Invalid, io.EOF), New(op, Invalid, io.EOF), true},
	{New(op, Invalid, io.EOF), New(op, Invalid, io.EOF), true},
	{New(op, Invalid, io.EOF), New(op, Invalid, io.EOF), true},
	{New(op, Invalid), New(op, Invalid, io.EOF), true},
	{New(op), New(op, Invalid, io.EOF), true},
	// Failure.
	{New(io.EOF), New(io.ErrClosedPipe), false},
	{New(op1), New(op2), false},
	{New(Invalid), New(Permission), false},
	// {E(op, Invalid, io.EOF, jane, path1), E(op, Invalid, io.EOF, john, path1), false},
	// Nested *Errors.
	{New(op1, New(op2)), New(op1, Str(New(op2).Error())), false},
}

func TestMatch(t *testing.T) {
	for _, test := range matchTests {
		matched := Match(test.err1, test.err2)
		if matched != test.matched {
			t.Errorf("Match(%q, %q)=%t; want %t", test.err1, test.err2, matched, test.matched)
		}
	}
}

type kindTest struct {
	err  error
	kind Kind
	want bool
}

var kindTests = []kindTest{
	// Non-Error errors.
	{nil, NotExist, false},
	{Str("not an *Error"), NotExist, false},

	// Basic comparisons.
	{New(NotExist), NotExist, true},
	{New(Exist), NotExist, false},
	{New("no kind"), NotExist, false},
	{New("no kind"), Other, false},

	// Nested *Error values.
	{New("Nesting", New(NotExist)), NotExist, true},
	{New("Nesting", New(Exist)), NotExist, false},
	{New("Nesting", New("no kind")), NotExist, false},
	{New("Nesting", New("no kind")), Other, false},
}

func TestKind(t *testing.T) {
	for _, test := range kindTests {
		got := Is(test.kind, test.err)
		if got != test.want {
			t.Errorf("Is(%q, %q)=%t; want %t", test.kind, test.err, got, test.want)
		}
	}
}
