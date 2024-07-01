package main

import (
	"errors"
	"math/rand/v2"
)

// simulateError returns an error about 30% of the time.
func simulateError() error {
	if rand.IntN(10) < 3 {
		return errors.New("unexpected error occurred")
	}
	return nil
}
