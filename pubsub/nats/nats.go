package pubsubnats

import (
	"github.com/go-kit/kit/log"
	nats "github.com/nats-io/nats.go"
	"github.com/go-kit/kit/log/level"
)

func WithDefaultConnectOptions(natsClientName string, logger log.Logger) ([]nats.Option) {
	var options []nats.Option
	options = append(options, nats.ReconnectHandler(func(c *nats.Conn) {
		level.Info(logger).Log("natsclient", natsClientName, "handler", "ReconnectHandler", "url", c.ConnectedUrl())
	}))
	options = append(options, nats.DisconnectHandler(func(c *nats.Conn) {
		level.Error(logger).Log("natsclient", natsClientName, "handler", "DisconnectHandler")
	}))
	options = append(options, nats.ClosedHandler(func(c *nats.Conn) {
		level.Error(logger).Log("natsclient", natsClientName, "handler", "ClosedHandler")
	}))
	options = append(options, nats.MaxReconnects(-1))
	options = append(options, nats.Name(natsClientName))

	return options
}
