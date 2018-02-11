package errors

type failableFunc func() error

// Monad is a container for error type which could be either nil or not
type Monad struct {
	err    error
	defers []func()
}

// Maybe returns a monad
func Maybe(fn failableFunc) Monad {
	return Monad{}.wrap(fn())
}

// Maybe returns a monad
func (e Monad) Maybe(fn failableFunc) Monad {
	if e.err != nil {
		return e
	}
	return e.wrap(fn())
}

// Err returns an error or nil
func (e Monad) Err() error {
	e.resolveDefers()
	if e.err != nil {
		return e.err
	}
	return nil
}

// Defer adds the defer func
func (e Monad) Defer(fn func()) Monad {
	if e.err != nil {
		return e
	}
	return Monad{e.err, append(e.defers, fn)}
}

func (e Monad) resolveDefers() {
	for _, fn := range e.defers {
		fn()
	}
}

func (e Monad) wrap(err error) Monad {
	return Monad{err, e.defers}
}
