package errors

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

type jsonMarshal struct {
	ID  Kind   `json:"id"`
	Op  Op     `json:"op,omitempty"`
	Msg string `json:"msg,omitempty"`
}

// MarshalJSON converts the err object to the JSON representation
func (e Error) MarshalJSON() ([]byte, error) {
	var err jsonMarshal
	if e.Kind > Other {
		err.ID = e.Kind
	}
	if e.Op != "" {
		err.Op = e.Op
	}
	if e.Err != nil {
		err.Msg = e.Error()
	}
	return json.Marshal(err)
}

// UnmarshalJSON deserializes JSON back to Error struct
func (e *Error) UnmarshalJSON(data []byte) error {
	var err jsonMarshal
	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}
	separatorLength := len(Separator)
	var trim int
	if err.Op != "" {
		e.Op = err.Op
		trim += len(e.Op) + separatorLength
	}
	if err.ID > Other {
		e.Kind = err.ID
		trim += len(e.Kind.String()) + separatorLength
	}
	if err.Msg != "" {
		e.Err = Str(err.Msg[trim:])
	}
	return nil
}

// MarshalJSON converts the err object to the JSON representation
func (e errorString) MarshalJSON() ([]byte, error) {
	var err jsonMarshal
	if e.s != "" {
		err.Msg = e.s
	}
	return json.Marshal(err)
}

// UnmarshalJSON deserializes JSON back to Error struct
func (e *errorString) UnmarshalJSON(data []byte) error {
	var err jsonMarshal
	if err := json.Unmarshal(data, &err); err != nil {
		return err
	}
	e.s = err.Msg
	return nil
}

func appendString(b []byte, str string) []byte {
	var tmp [16]byte // For use by PutUvarint.
	N := binary.PutUvarint(tmp[:], uint64(len(str)))
	b = append(b, tmp[:N]...)
	b = append(b, str...)
	return b
}

// getBytes unmarshals the byte slice at b (uvarint count followed by bytes)
// and returns the slice followed by the remaining bytes.
// If there is insufficient data, both return values will be nil.
func getBytes(b []byte) (data, remaining []byte) {
	u, N := binary.Uvarint(b)
	if len(b) < N+int(u) {
		log.Printf("Unmarshal error: bad encoding")
		return nil, nil
	}
	if N == 0 {
		log.Printf("Unmarshal error: bad encoding")
		return nil, b
	}
	return b[N : N+int(u)], b[N+int(u):]
}
