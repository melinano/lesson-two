package dev07

import (
	"testing"
	"time"
)

func TestOr(t *testing.T) {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	t.Run("returns nil for no channels", func(t *testing.T) {
		result := or()
		if result != nil {
			t.Error("Expected nil, but got a non-nil channel.")
		}
	})

	t.Run("returns after shortest duration for multiple channels", func(t *testing.T) {
		start := time.Now()
		<-or(
			sig(1*time.Second),
			sig(2*time.Second),
			sig(3*time.Second),
		)
		duration := time.Since(start)
		if duration < 1*time.Second || duration > 2*time.Second {
			t.Errorf("Expected to return after 1 second, but returned after %v.", duration)
		}
	})
}
