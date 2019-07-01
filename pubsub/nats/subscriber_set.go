package pubsubnats

import (
	"fmt"
	"io"
	"strings"

	nats "github.com/nats-io/nats.go"
)

type SubscriptionSet struct {
	io.Closer
	Err           error
	subscriptions []*nats.Subscription
}

type SubscriberFunc func() (*nats.Subscription, error)

func RegisterSubscribers(funcs ...SubscriberFunc) SubscriptionSet {
	return SubscriptionSet{}.Register(funcs...)
}

func (registry SubscriptionSet) Register(funcs ...SubscriberFunc) SubscriptionSet {
	for _, f := range funcs {
		uSub, err := f()
		if err != nil {
			registry.Err = err
			break
		}
		registry.subscriptions = append(registry.subscriptions, uSub)
	}
	return registry
}

func (registry SubscriptionSet) Close() error {
	var errstrings []string
	for _, sub := range registry.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			errstrings = append(errstrings, err.Error())
		}
	}
	return fmt.Errorf(strings.Join(errstrings, "\n"))
}
