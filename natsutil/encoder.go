package natsutil

import (
	"context"
	"encoding/json"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"

	"github.com/etherlabsio/errors"
	"github.com/go-kit/kit/log"
	natstransport "github.com/go-kit/kit/transport/nats"
	nats "github.com/nats-io/go-nats"
)

// JSONErrorEncoder is a nats RPC JSON reply error encoder
func JSONErrorEncoder(l log.Logger) natstransport.ErrorEncoder {
	return func(ctx context.Context, err error, reply string, nc *nats.Conn) {
		if err == nil {
			panic("nats.JSONErrEncoder received nil error")
		}
		e := errors.Serializable(err)
		b, err := json.Marshal(map[string]interface{}{
			"error": e,
		})
		if err != nil {
			level.Error(l).Log("msg", "marshal nats error failure", "err", err)
			return
		}
		if err := nc.Publish(reply, b); err != nil {
			level.Error(l).Log("msg", "nats error reply publish failure", "err", err, "reply", reply)
		}
	}
}

// JSONResponseEncoder is a EncodeResponseFunc that serializes the response as a
// JSON object to the subscriber reply. Many JSON-over services can use it as
// a sensible default.
func JSONResponseEncoder(encodeErr natstransport.ErrorEncoder) natstransport.EncodeResponseFunc {
	return func(ctx context.Context, reply string, nc *nats.Conn, response interface{}) error {
		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			encodeErr(ctx, f.Failed(), reply, nc)
			return nil
		}
		b, err := json.Marshal(response)
		if err != nil {
			return errors.WithMessagef(err, "failed to marshal response for reply %s", reply)
		}
		return nc.Publish(reply, b)
	}
}
