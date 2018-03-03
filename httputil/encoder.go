package httputil

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/etherlabs/pkg/errors"
)

var errYouAreDoingItWrong = errors.New("programmer error")

// DefaultErrorEncoder takes in a status coder and returns an HTTP error encoder
func DefaultErrorEncoder(statusCoder func(err error) int) func(_ context.Context, err error, w http.ResponseWriter) {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			err = errors.WithMessage(errYouAreDoingItWrong, "encodeError received nil error")
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCoder(err))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
	}
}
