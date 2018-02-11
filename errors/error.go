// Package errors defines the error handling used by all Ether backend systems.
// Function definition inspiration derived from github.com/pkg/errors
// Design inpiration derived from github.com/upspin/errors
// Thank you Rob Pike :P
package errors

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Op is the operation being performed, usually the name of the method
	// being invoked (Get, Put, etc.). It should not contain an at sign @.
	Op Op
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind
	// The underlying error that triggered this one, if any.
	Err error

	// Used along with the wrap function to bind more context to the error
	// when it is propagated back up the call chain
	msg

	*stack
}

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == 0 && e.Err == nil
}

var (
	_ error = (*Error)(nil)
)

type debug bool

const Trace debug = true

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ": "

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind int

func (k Kind) String() string {
	return strconv.Itoa(int(k))
}

type msg string

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Other      Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                // Invalid operation for this type of item.
	Permission             // Permission denied.
	IO                     // External I/O error such as network failure.
	Exist                  // Item already exists.
	NotExist               // Item does not exist.
	Internal               // Internal error or inconsistency.
)

// New builds an error value from its arguments.
// There must be at least one argument or New panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	errors.Op
//		The operation being performed, usually the method
//		being invoked (Get, Put, etc.).
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.Str.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func New(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case string:
			e.Err = Str(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		case msg:
			e.msg = arg
		case debug:
			e.stack = callers()
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}

	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}

	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind twice.
	if prev.Kind == e.Kind {
		prev.Kind = Other
	}

	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind, prev.Kind = prev.Kind, Other
	}

	// If this error has Op unset, pull up the inner one.
	if e.Op == "" {
		e.Op, prev.Op = prev.Op, ""
	}

	e.stack, prev.stack = prev.stack, nil

	return e
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, Separator)
		b.WriteString(string(e.Op))
	}
	if e.Kind != 0 {
		pad(b, Separator)
		b.WriteString(e.Kind.String())
	}
	if e.msg != "" {
		pad(b, Separator)
		b.WriteString(string(e.msg))
	}
	if e.Err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, Separator)
			b.WriteString(e.Err.Error())
		}
	}
	if e.stack != nil {
		pad(b, Separator)
		b.WriteString(fmt.Sprintf("%+v", e.stack.StackTrace()))
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
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}

// Match compares its two error arguments. It can be used to check
// for expected errors in tests. Both arguments must have underlying
// type *Error or Match will return false. Otherwise it returns true
// iff every non-zero element of the first error is equal to the
// corresponding element of the second.
// If the Err field is a *Error, Match recurs on that field;
// otherwise it compares the strings returned by the Error methods.
// Elements that are in the second argument but not present in
// the first are ignored.
//
// For example,
//	Match(errors.E(errors.Permission), err)
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != Other && e2.Kind != e1.Kind {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}

// Is reports whether err is an *Error of the given Kind.
// If err is nil then Is returns false.
func Is(kind Kind, err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.Kind != Other {
		return e.Kind == kind
	}
	if e.Err != nil {
		return Is(kind, e.Err)
	}
	return false
}
