package logger

import (
	"sync"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	c := newConfig(LevelDebug)
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for i := 100; i > 0; i-- {
			c.setLevel(LevelWarning)
		}

		wg.Done()
	}()

	go func() {
		for i := 100; i > 0; i-- {
			c.setLevel(LevelError)
		}

		wg.Done()
	}()

	go func() {
		for i := 100; i > 0; i-- {
			c.getLevel()
		}

		wg.Done()
	}()

	wg.Wait()
}
