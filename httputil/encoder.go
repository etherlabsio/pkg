package httputil

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/etherlabsio/errors"
)

var errYouAreDoingItWrong = errors.New("programmer error")

type SerializerFunc func(w http.ResponseWriter, v interface{}) error

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

func JSONSerializer(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}
