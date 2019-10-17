package httputil

import (
	"net/http"

	"github.com/etherlabsio/errors"
)

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range or equal to 202 Accepted.
// API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
//
// The error type will be *TwoFactorAuthError for two-factor authentication errors.
func CheckResponse(r *http.Response) error {
	success := func(statusCode int) bool {
		return 200 <= statusCode && statusCode < 500
	}(r.StatusCode)
	if success {
		return nil
	}
	return errors.New("internal server error with http status: "+r.Status, errors.Internal)
}

// CheckOKResponse checks if the status code is within 2XX range
func CheckOKResponse(r *http.Response) error {
	isOK := r.StatusCode >= 200 && r.StatusCode <= 299
	if !isOK {
		return errors.New("response error with http status: "+r.Status, errors.Internal)
	}
	return nil
}
