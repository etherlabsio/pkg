package errors

import (
	"encoding/json"
)

type jsonError struct {
	Code  Kind   `json:"code"`
	Msg   string `json:"message,omitempty"`
	Cause string `json:"cause,omitempty"`
}

// MarshalJSON converts the err object to the JSON representation
func (e *Error) MarshalJSON() ([]byte, error) {
	var err jsonError
	err.Code = e.kind
	err.Msg = e.msg
	if e.cause != nil {
		err.Cause = e.Cause().Error()
	}
	return json.Marshal(err)
}

// UnmarshalJSON deserializes JSON back to Error struct
func (e *Error) UnmarshalJSON(data []byte) error {
	var err jsonError
	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}
	e.kind = err.Code
	e.msg = err.Msg
	if cause := err.Cause; cause != "" {
		e.cause = Str(cause)
	}
	return nil
}

// MarshalJSON converts the err object to the JSON representation
func (e *stringerr) MarshalJSON() ([]byte, error) {
	var err jsonError
	err.Code = Internal
	err.Msg = e.msg
	return json.Marshal(err)
}
