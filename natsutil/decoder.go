package natsutil

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/etherlabsio/errors"
	natstransport "github.com/go-kit/kit/transport/nats"
	nats "github.com/nats-io/go-nats"
)

// Caution: Do not cargo cult this code. This is meant for specific use case based on our understanding of the message payload possibilities.
func DecodeNATSJSONEvent(decoderMap map[string]interface{}) natstransport.DecodeRequestFunc {
	return func(_ context.Context, msg *nats.Msg) (request interface{}, err error) {
		event, ok := decoderMap[msg.Subject]
		if !ok {
			return nil, errors.Errorf("decoder type for event %s undefined", msg.Subject)
		}

		v := reflect.New(reflect.TypeOf(event)).Interface()
		if err := json.Unmarshal(msg.Data, v); err != nil {
			return nil, errors.WithMessage(err, "event decoding failed for subject "+msg.Subject)
		}

		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Ptr {
			return val.Elem().Interface(), nil
		}
		return val.Interface(), nil
	}
}
