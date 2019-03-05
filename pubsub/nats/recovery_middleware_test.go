package pubsubnats

import (
	"testing"

	"github.com/etherlabsio/pkg/logutil"
	"github.com/nats-io/go-nats"
	"github.com/stretchr/testify/assert"
)

type mockRecoveryHandler struct{}

func (h mockRecoveryHandler) ServeMsg(nc *nats.Conn) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		panic("something went wrong")
	}
}

func TestRecoveryMiddleware_ServeMsg(t *testing.T) {
	var h Handler
	h = mockRecoveryHandler{}
	h = NewRecoveryMiddleware(logutil.NewServerLogger(true), h)

	assert := assert.New(t)

	assert.NotPanics(func() {
		msgHandler := h.ServeMsg(nil)
		msgHandler(&nats.Msg{
			Subject: "test_subject",
			Data:    []byte("test payload"),
		})
	})
}
