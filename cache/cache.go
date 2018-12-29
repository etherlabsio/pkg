package cache

import "time"

const (
	// NoExpiry specifies no expiry to the key for which the value is set
	NoExpiry time.Duration = 0
)

// Cache interface is a generic cache definition used for most common types of cacheing operations
type Cache interface {
	Set(key string, value interface{}, expiry time.Duration) bool
	Get(key string, marshallableValue interface{}) bool
	Delete(key string) bool
}
