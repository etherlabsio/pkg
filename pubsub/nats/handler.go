package pubsubnats

import "github.com/nats-io/go-nats"

// Handler serves messages for NATS
type Handler interface {
	ServeMsg(nc *nats.Conn) func(msg *nats.Msg)
}
