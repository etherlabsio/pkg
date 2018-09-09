package errors

type causer interface {
	Cause() error
}

type kinder interface {
	Kind() Kind
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

// KindOf returns the underlying kind type of the error, if possible.
// An error value has a kind if it implements the following
// interface:
//
//     type kinder interface {
//            Kind() Kind
//     }
//
// If the error does not implement Kind, the original error will
// be returned. If the error is nil, Internal will be returned without further
// investigation.
func KindOf(err error) Kind {
	for err != nil {
		kind, ok := err.(kinder)
		if ok {
			return kind.Kind()
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return Internal
}

func errWithKind(err error) *withKind {
	for err != nil {
		e, ok := err.(*withKind)
		if ok {
			return e
		}
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return &withKind{
		msg:  "Internal error or inconsistency",
		kind: Internal,
	}
}
