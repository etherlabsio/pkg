package pubsub

import (
	"context"
)

// Publisher is a minimal interface for publishing messages to a pool of
// subscribers. Publishers are probably (but not necessarily) sending to a
// message bus.
//
// Most paramaterization of the publisher (topology restrictions like a topic,
// exchange, or specific message type; queue or buffer sizes; etc.) should be
// done in the concrete constructor.
type Publisher interface {
	// Publish a single message, described by an io.Reader, to the given key.
	//
	// CHANGE(chrisprijic): added context, so transport can utilize it.
	Publish(ctx context.Context, key string, msg interface{}) error
}
