package redis

import (
	"bytes"
	"encoding/gob"
	"log"
	"testing"

	"github.com/alicebob/miniredis"
)

func TestCache_CheckUmarshalling(t *testing.T) {
	// miniredis is an in-mememory implementation of Redis protocol.
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	const ns = "test"
	const key = "key"

	client := NewClient(
		Addresses(s.Addr()),
	)

	c := NewCachev2(client, Namespace(ns))

	t.Run("get and set string", func(t *testing.T) {
		input := "something"

		var expected string
		c.Set(key, input, 0)
		c.Get(key, &expected)
		t.Log("found: ", expected)
		if expected != input {
			t.Errorf("redis.cacheV2.get() got = %v, wanted = %v", expected, input)
		}
	})

	t.Run("get and set binary", func(t *testing.T) {
		type binaryType struct {
			Name  string
			Value int
		}

		input := binaryType{"karthik", 27}

		var b bytes.Buffer
		err := gob.NewEncoder(&b).Encode(input)
		if err != nil {
			log.Fatalf("input marshalling error: %v", err)
		}

		var expected binaryType
		c.Set(key, input, 0)
		c.Get(key, &expected)
		if (expected.Name != input.Name) || (expected.Value != input.Value) {
			t.Errorf("redis.cacheV2.get() got = %v, want %v", expected, input)
		}
	})

	t.Run("get and set boolean", func(t *testing.T) {
		const input = true

		var expected bool
		c.Set(key, input, 0)
		c.Get(key, &expected)
		if expected != input {
			t.Errorf("redis.cacheV2.get() got = %v, want %v", expected, input)
		}
	})
}
