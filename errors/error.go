package errors

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

const separator = ":: "
const opSeparator = ": "

var _ error = (*Error)(nil)

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind int

func (k Kind) String() string {
	return strconv.Itoa(int(k))
}

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Internal     Kind = iota // Internal error or inconsistency.
	Invalid                  // Invalid operation for this type of item.
	Permission               // Permission denied.
	IO                       // External I/O error such as network failure.
	AlreadyExist             // Item already exists.
	NotExist                 // Item does not exist.
)

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked (Get, Put, etc.). It should not contain an at sign @.
	op Op
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	*withKind
}

// New builds an error value from its arguments.
// There must be at least one argument or New panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	errors.op
//		The operation being performed, usually the method
//		being invoked (Get, Put, etc.).
//	errors.kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Internal, we set it to the Kind of
// the underlying error.
//
func New(msg string, args ...interface{}) error {
	e := &Error{
		withKind: &withKind{
			msg:  msg,
			kind: Internal,
		},
	}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.op = arg
		case Kind:
			e.kind = arg
		case error:
			e.cause = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	prev, ok := e.cause.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind twice.
	if prev.kind == e.kind {
		return e
	}

	// If this error has Kind unset or Other, pull up the inner one.
	if e.kind == Internal {
		e.kind = prev.kind
	}
	return e
}

// Cause implements the causer interface
func (e *Error) Cause() error {
	return e.withKind.Cause()
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	b.WriteString(string(e.msg))
	if e.op != "" {
		pad(b, separator)
		b.WriteString(string(e.op))
	}
	if e.cause != nil {
		pad(b, separator)
		b.WriteString(e.cause.Error())
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// Recreate the errors.New functionality of the standard Go errors package
// so we can create simple text errors when needed.

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	if text == "" {
		_, file, line, _ := runtime.Caller(1)
		return Errorf("errors.E: bad call from %s:%d", file, line)
	}
	return &fundamental{text}
}

// fundamental is a trivial implementation of error.
type fundamental struct {
	msg string
}

func (e *fundamental) Error() string {
	return e.msg
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &fundamental{msg: fmt.Sprintf(format, args...)}
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}
