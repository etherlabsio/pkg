package errorutil

import (
	"context"

	"github.com/etherlabsio/errors"
	"github.com/go-kit/kit/endpoint"
)

type errResponse struct {
	Err *errors.Error `json:"error,omitempty"`
}

func (e errResponse) Failed() error {
	if e.Err != nil {
		return e.Err
	}
	return nil
}

// UnwrapMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func UnwrapMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			response, err = next(ctx, request)
			if err != nil {
				return nil, errors.Unwrap(err)
			}
			if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
				// It should be logged already on the service-level logger middleware
				err := errors.Unwrap(err)
				return errResponse{Err: err}, nil
			}
			return response, nil
		}
	}
}
