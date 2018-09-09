package errors

type causer interface {
	Cause() error
}

type kinder interface {
	Kind() Kind
}
