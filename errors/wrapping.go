package errors

import (
	"fmt"

	errors "github.com/pkg/errors"
)

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	return errors.WithStack(err)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is call, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

type withOp struct {
	op    Op
	cause error
}

func (err *withOp) Error() string {
	return string(err.op) + separator + err.cause.Error()
}

func (err *withOp) Cause() error {
	return err.cause
}

// WithOp returns an error annotating err with a hint of the operation name
// at the point WithOp is called. If err is nil, WithOp returns nil.
func WithOp(err error, op Op) error {
	if err == nil {
		return nil
	}
	return &withOp{
		cause: err,
		op:    op,
	}
}

// WithOpf returns an error annotating err with a hint of the operation name
// at the point WithOpf is call, and the format specifier.
// If err is nil, WithOpf returns nil.
func WithOpf(err error, op Op, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withOp{
		cause: errors.WithMessage(err, fmt.Sprintf(format, args...)),
		op:    op,
	}
}

type withKind struct {
	kind  Kind
	msg   string
	cause error
}

func (e *withKind) isZero() bool {
	return e.msg == "" && e.kind == Internal && e.cause == nil
}

func (e *withKind) Error() string {
	return e.msg + separator + e.cause.Error()
}

func (e *withKind) Cause() error {
	if e.isZero() {
		return nil
	}
	return e.cause
}

func (e *withKind) Kind() Kind {
	return e.kind
}

// WithKind returns an error annotating err with the service specific kind of err
// at the point WithKind is called. If err is nil, WithKind returns nil.
func WithKind(err error, kind Kind, msg string) error {
	if err == nil {
		return nil
	}
	return &withKind{
		cause: err,
		msg:   msg,
		kind:  kind,
	}
}

// WithKindf returns an error annotating err with the service specific kind of err
// at the point WithKindf is call, and the format specifier.
// If err is nil, WithKindf returns nil.
func WithKindf(err error, kind Kind, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withKind{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		kind:  kind,
	}
}
