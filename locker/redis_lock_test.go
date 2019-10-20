package locker

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-redis/redis"

	"github.com/alicebob/miniredis"
)

// Additional test cases can be found at https://github.com/bsm/redislock/blob/master/redislock_test.go

func setup() *miniredis.Miniredis {
	// miniredis is an in-mememory implementation of Redis protocol.
	s, err := miniredis.Run()
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func TestRedisLocker_LockAndUnlock(t *testing.T) {
	db := setup()
	defer db.Close()

	c := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    db.Addr(),
	})

	t.Run("first lock acquire wins, others fail", func(t *testing.T) {
		locker := NewRedisLocker(c)
		l1, err := locker.Lock(context.Background(), "test1", WithTTL(5*time.Second))
		assert.Nil(t, err)
		defer l1.Unlock()

		_, err = locker.Lock(context.Background(), "test1", WithTTL(2*time.Second))
		assert.Equal(t, ErrNotObtained, err)
	})

	t.Run("multiple locks can be acquired on separate keys", func(t *testing.T) {
		locker := NewRedisLocker(c)
		l1, err := locker.Lock(context.Background(), "test3", WithTTL(5*time.Second))
		assert.Nil(t, err)
		err = l1.Unlock()
		assert.Nil(t, err)

		l2, err := locker.Lock(context.Background(), "test4", WithTTL(2*time.Second))
		assert.Nil(t, err)
		err = l2.Unlock()
		assert.Nil(t, err)
	})

	t.Run("should release lock on ttl expiry", func(t *testing.T) {
		locker := NewRedisLocker(c)

		l1, err := locker.Lock(context.Background(), "test5", WithTTL(time.Millisecond))
		assert.Nil(t, err)

		time.Sleep(5 * time.Millisecond)
		assert.Equal(t, ErrNotHeld, l1.Unlock())
	})
}
