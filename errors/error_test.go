package errors

import (
	"testing"
)

func TestKind_String(t *testing.T) {
	kind := Permission
	const expected = "2"
	if kind.String() != expected {
		t.Fatalf("expected %s, got %s", expected, kind.String())
	}
}

func TestDoesNotChangePreviousError(t *testing.T) {
	err := New("new error", Permission)
	err2 := New("wrapped context", Op("op"), err)

	expected := "op" + separator + "wrapped context" + separator + "new error"
	if err2.Error() != expected {
		t.Fatalf("Expected %q, got %q", expected, err2)
	}
	kind := err.(kinder).Kind()
	if kind != Permission {
		t.Fatalf("Expected kind %v, got %v", Permission, kind)
	}
	l3op := Op("someFunc")
	err3 := New("level3 err", l3op, Internal, err2)
	e := err3.(*withOp)
	if e.op != l3op {
		t.Fatalf("Expected op %v, got %v", l3op, e.op)
	}
}

func TestError_Error(t *testing.T) {
	err1 := New("new error", Permission)
	err2 := New("wrapped context", Op("I will NOT modify err"), err1)
	l3op := Op("someFunc")
	err3 := New("level3 err", l3op, Internal, err2)
	type fields struct {
		op    Op
		Error *Error
	}
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			"level1 Error()",
			err1,
			"new error",
		},
		{
			"level2 Error()",
			err2,
			"I will NOT modify err: wrapped context: new error",
		},
		{
			"level3 Error()",
			err3,
			"someFunc: level3 err: I will NOT modify err: wrapped context: new error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStr(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"correct call",
			args{text: "test"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Str(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("Str() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stringerr_Error(t *testing.T) {
	type fields struct {
		msg string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "new stringerr error",
			fields: fields{"test new error"},
			want:   "test new error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &stringerr{
				msg: tt.fields.msg,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("stringerr.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
