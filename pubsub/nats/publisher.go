package pubsubnats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/etherlabsio/errors"
	"github.com/go-kit/kit/log/level"
	natstransport "github.com/go-kit/kit/transport/nats"

	"github.com/go-kit/kit/log"
	"github.com/nats-io/go-nats"
)

// EncodeRequestFunc encodes the passed request object into the NATS request
// object. It's designed to be used in NATS publishers, for publisher-side
// endpoints. One straightforward EncodeRequestFunc could something that JSON
// encodes the object directly to the request payload.
type EncodeRequestFunc natstransport.EncodeRequestFunc

// RequestFunc may take information from a publisher request and put it into a
// request context. In Subscribers, RequestFuncs are executed prior to invoking the
// endpoint.
type RequestFunc natstransport.RequestFunc

// Publisher wraps a URL and provides a method that implements endpoint.Endpoint.
type Publisher struct {
	publisher *nats.Conn
	enc       EncodeRequestFunc
	before    []RequestFunc
	after     []RequestFunc
	logger    log.Logger
	timeout   time.Duration
}

// NewPublisher constructs a usable Publisher for a single remote method.
func NewPublisher(
	publisher *nats.Conn,
	options ...PublisherOption,
) *Publisher {
	p := &Publisher{
		publisher: publisher,
		enc:       EncodeJSONRequest,
		logger:    log.NewNopLogger(),
		timeout:   10 * time.Second,
	}
	for _, option := range options {
		option(p)
	}
	return p
}

// PublisherOption sets an optional parameter for clients.
type PublisherOption func(*Publisher)

// PublisherBefore sets the RequestFuncs that are applied to the outgoing NATS
// request before it's invoked.
func PublisherBefore(before ...RequestFunc) PublisherOption {
	return func(p *Publisher) { p.before = append(p.before, before...) }
}

// PublisherTimeout sets the available timeout for NATS request.
func PublisherTimeout(timeout time.Duration) PublisherOption {
	return func(p *Publisher) { p.timeout = timeout }
}

func PublisherLogger(l log.Logger) PublisherOption {
	return func(p *Publisher) { p.logger = log.With(l, "component", "messaging_publisher") }
}

func PublisherVerbose() PublisherOption {
	return func(p *Publisher) {
		before := func(ctx context.Context, msg *nats.Msg) context.Context {
			level.Info(p.logger).Log(
				"status", "pending",
				"topic", msg.Subject,
				"payload", string(msg.Data),
			)
			return ctx
		}

		p.before = append(p.before, before)

		after := func(ctx context.Context, msg *nats.Msg) context.Context {
			level.Info(p.logger).Log(
				"status", "completed",
				"topic", msg.Subject,
				"payload", string(msg.Data),
			)
			return ctx
		}

		p.after = append(p.after, after)
	}
}

// Publish returns a usable endpoint that invokes the remote endpoint.
func (p Publisher) Publish(ctx context.Context, subject string, e interface{}) error {
	const op = "nats.Publish"
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	msg := nats.Msg{Subject: subject}

	if err := p.enc(ctx, &msg, e); err != nil {
		level.Error(p.logger).Log(
			"topic", msg.Subject,
			"status", "failed",
			"err", err,
		)
		return errors.Errorf("%s: encoder failure for topic %s and payload %+v", op, subject, e)
	}

	for _, f := range p.before {
		ctx = f(ctx, &msg)
	}

	err := p.publisher.Publish(msg.Subject, msg.Data)
	if err != nil {
		level.Error(p.logger).Log(
			"topic", msg.Subject,
			"status", "failed",
			"err", err,
		)
		return errors.Errorf("%s: publish failure for topic %s", op, subject)
	}

	for _, f := range p.after {
		ctx = f(ctx, &msg)
	}

	return nil
}

// EncodeJSONRequest is an EncodeRequestFunc that serializes the request as a
// JSON object to the Data of the Msg. Many JSON-over-NATS services can use it as
// a sensible default.
func EncodeJSONRequest(_ context.Context, msg *nats.Msg, request interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}
	msg.Data = b
	return nil
}
