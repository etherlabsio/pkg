package errors

import "testing"

func Test_Monad(t *testing.T) {
	wantDefer := "defer1defer2"
	gotDefer := ""
	want := "something is wrong"
	maybe := Maybe(func() error {
		return nil
	}).Defer(func() {
		gotDefer += "defer1"
	})
	maybe = maybe.Maybe(func() error {
		return nil
	}).Defer(func() {
		gotDefer += "defer2"
	})
	maybe = maybe.Maybe(func() error {
		return New(want)
	})
	maybe = maybe.Maybe(func() error {
		return Str("nothing")
	}).Defer(func() {
		gotDefer += "defer3"
	})
	err := maybe.Err()
	if err == nil {
		t.Errorf("Test_Maybe(): expected error; got nil")
	}
	if err.Error() != want {
		t.Errorf("Test_Maybe(): expected err string %s; got %s", want, err.Error())
	}
	if gotDefer != wantDefer {
		t.Errorf("Tet_Defer(): expected %s; got %s", wantDefer, gotDefer)
	}
}
