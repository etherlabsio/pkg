package pubsubnats

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

// RecoveryMiddleware intercepts the nats messages and performs recovery on panic
type RecoveryMiddleware struct {
	logger log.Logger
	next   Handler
}

// NewRecoveryMiddleware returns as RecoveryMiddleware handler
func NewRecoveryMiddleware(l log.Logger, h Handler) Handler {
	return RecoveryMiddleware{l, h}
}

// ServeMsg wraps the serveMsg handler with the stacktrace
func (mw RecoveryMiddleware) ServeMsg(nc *nats.Conn) func(msg *nats.Msg) {
	handler := mw.next.ServeMsg(nc)
	return func(msg *nats.Msg) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err := errors.New(fmt.Sprintf("Panic: %+v", rvr))
				err = errors.WithStack(err)
				mw.logger.Log(
					"err", fmt.Sprintf("%+v", err),
					"subject", msg.Subject,
				)
			}
		}()

		handler(msg)
	}
}
