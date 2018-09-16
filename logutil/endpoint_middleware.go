package logutil

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

// LoggingMiddleware returns an endpoint middleware that logs the
// duration of each invocation, and the resulting error, if any.
func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
					// It should be logged already on the service-level logger middleware
					return
				}
				WithError(logger, err).Log("component", "endpoint", "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
