package errors

import "fmt"

type formatFunc func() string

// Wrap adds context to the err when it is propagated back up the call chain
// When you call Error() method, you should see the additional information
func Wrap(err error, s string) error {
	return wrap(err, func() string { return s })
}

// Wrapf adds context to the err when it is propagated back up the call chain
// When you call Error() method, you should see the additional information.
// Instead of passing a string, you can pass formattable arguments
func Wrapf(err error, format string, args ...interface{}) error {
	return wrap(err, func() string {
		return fmt.Sprintf(format, args...)
	})
}

// WithMessage is currently just a shim for the func with same name in github.com/pkg/errors
func WithMessage(err error, s string) error {
	return Wrap(err, s)
}

func wrap(err error, s formatFunc) error {
	if err == nil {
		return nil
	}
	str := s()
	err = New(err, msg(str))
	return err
}
