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

// JSONErrorEncoder takes in a status coder and returns an HTTP error encoder
func JSONErrorEncoder(statusCoder func(err error) int) httptransport.ErrorEncoder {
	return func(_ context.Context, err error, w http.ResponseWriter) {
		if err == nil {
			err = errors.WithMessage(errYouAreDoingItWrong, "encodeError received nil error")
			panic(err)
		}
		w.WriteHeader(statusCoder(err))
		e := errors.Serializable(err)
		JSONSerializer(w, map[string]interface{}{
			"error": e,
		})
	}
}

// EncodeJSONResponse encodes a response using the appropriate serializer function
func EncodeJSONResponse(encodeErr httptransport.ErrorEncoder) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
			encodeErr(ctx, f.Failed(), w)
			return nil
		}

		if headerer, ok := response.(httptransport.Headerer); ok {
			for k, values := range headerer.Headers() {
				for _, v := range values {
					w.Header().Add(k, v)
				}
			}
		}
		code := http.StatusOK
		if sc, ok := response.(httptransport.StatusCoder); ok {
			code = sc.StatusCode()
		}
		w.WriteHeader(code)
		if code == http.StatusNoContent {
			return nil
		}
		return JSONSerializer(w, response)
	}
}

// JSONSerializer returns a encodes the data to a JSON response
func JSONSerializer(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(v)
}
