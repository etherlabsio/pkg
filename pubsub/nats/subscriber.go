package pubsubnats

import (
	"context"

	natstransport "github.com/go-kit/kit/transport/nats"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/nats-io/go-nats"
)

// DecodeRequestFunc extracts a user-domain request object from a publisher
// request object. It's designed to be used in NATS subscribers, for subscriber-side
// endpoints. One straightforward DecodeRequestFunc could be something that
// JSON decodes from the request body to the concrete response type.
type DecodeRequestFunc natstransport.DecodeRequestFunc

// Subscriber wraps an endpoint and provides nats.MsgHandler.
type Subscriber struct {
	e      endpoint.Endpoint
	dec    DecodeRequestFunc
	before []RequestFunc
	logger log.Logger
}

// NewSubscriber constructs a new subscriber, which provides nats.MsgHandler and wraps
// the provided endpoint.
func NewSubscriber(
	e endpoint.Endpoint,
	dec DecodeRequestFunc,
	options ...SubscriberOption,
) *Subscriber {
	s := &Subscriber{
		e:      e,
		dec:    dec,
		logger: log.NewNopLogger(),
	}
	for _, option := range options {
		option(s)
	}
	return s
}

// SubscriberOption sets an optional parameter for subscribers.
type SubscriberOption func(*Subscriber)

// SubscriberBefore functions are executed on the publisher request object before the
// request is decoded.
func SubscriberBefore(before ...RequestFunc) SubscriberOption {
	return func(s *Subscriber) { s.before = append(s.before, before...) }
}

// SubscriberErrorLogger is used to log non-terminal errors. By default, no errors
// are logged. This is intended as a diagnostic measure. Finer-grained control
// of error handling, including logging in more detail, should be performed in a
// custom SubscriberErrorEncoder which has access to the context.
func SubscriberErrorLogger(logger log.Logger) SubscriberOption {
	return func(s *Subscriber) { s.logger = log.With(level.Error(logger), "component", "messaging_subscriber") }
}

// Serve provides nats.MsgHandler.
func (s Subscriber) ServeMsg(nc *nats.Conn) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, f := range s.before {
			ctx = f(ctx, msg)
		}

		request, err := s.dec(ctx, msg)
		if err != nil {
			s.logger.Log(
				"msg", "error decoding nats msg",
				"subject", msg.Subject,
				"err", err,
			)
			errHandler(ctx, err, msg, nc)
			return
		}

		_, err = s.e(ctx, request)
		if err != nil {
			s.logger.Log(
				"msg", "endpoint error for nats msg",
				"subject", msg.Subject,
				"err", err,
			)
			errHandler(ctx, err, msg, nc)
			return
		}
	}
}

// TODO: Add a way to handle errors
func errHandler(ctx context.Context, err error, msg *nats.Msg, nc *nats.Conn) {}

// NopRequestDecoder is a DecodeRequestFunc that can be used for requests that do not
// need to be decoded, and simply returns nil, nil.
func NopRequestDecoder(_ context.Context, _ *nats.Msg) (interface{}, error) {
	return nil, nil
}
