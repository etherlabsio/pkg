package httputil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/etherlabsio/errors"
)

var errYouAreDoingItWrong = errors.New("programmer error")

// SerializerFunc serializes a specific response
type SerializerFunc func(w http.ResponseWriter, v interface{}) error

// ErrorEncoderFunc encodes a response if it contains an error
type ErrorEncoderFunc func(_ context.Context, err error, w http.ResponseWriter)

// DefaultErrorEncoder takes in a status coder and returns an HTTP error encoder
func DefaultErrorEncoder(encode SerializerFunc, statusCoder func(err error) int) func(_ context.Context, err error, w http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			err = errors.WithMessage(errYouAreDoingItWrong, "encodeError received nil error")
		}
		w.WriteHeader(statusCoder(err))
		encode(w, map[string]interface{}{
			"error": err,
		})
	}
}

// ResponseEncoder encodes a response using the appropriate serializer function
func ResponseEncoder(s SerializerFunc, e ErrorEncoderFunc) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			e(ctx, f.Failed(), w)
			return nil
		}
		return s(w, response)
	}
}

// JSONSerializer returns a encodes the data to a JSON response
func JSONSerializer(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}
