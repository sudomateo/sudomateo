package main

import (
	"context"
	"testing"
	"testing/synctest"
	"time"
)

func TestGOMAXPROCSWatcher_Stop(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Start the watcher with a context that never cancels so we can test that
		// calling [GOMAXPROCSWatcher.Stop] stops the watcher.
		watcher := NewGOMAXPROCSWatcher()
		ch := watcher.Watch(context.Background(), time.Nanosecond)

		// Check that the watcher channel is still open.
		synctest.Wait()
		select {
		case _, ok := <-ch:
			if !ok {
				t.Fatal("watcher channel is closed on a running watcher")
			}
		default:
		}

		// Stop the watcher.
		watcher.Stop()
		synctest.Wait()

		// Confirm that the watcher channel is closed.
		if _, ok := <-ch; ok {
			t.Fatalf("watcher channel is open on a stopped watcher")
		}
	})
}

func TestGOMAXPROCSWatcher_Context(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		timeout := 5 * time.Nanosecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Start the watcher with a context that cancels after a timeout so we can test
		// that a cancelled context stops the watcher.
		watcher := NewGOMAXPROCSWatcher()
		ch := watcher.Watch(ctx, time.Nanosecond)

		// Wait just before the context is cancelled.
		time.Sleep(timeout - time.Nanosecond)
		synctest.Wait()

		// Confirm the context is still valid.
		if err := ctx.Err(); err != nil {
			t.Fatalf("context cancelled before timeout: %v", err)
		}

		// Check that the watcher channel is still open.
		select {
		case _, ok := <-ch:
			if !ok {
				t.Fatal("watcher channel is closed on a running watcher")
			}
		default:
		}

		// Wait the rest of the time needed to cancel the context.
		time.Sleep(time.Nanosecond)
		synctest.Wait()

		// Confirm the context cancelled successfully.
		if err := ctx.Err(); err != context.DeadlineExceeded {
			t.Fatalf("context cancelled with unexpected error: %v", err)
		}

		// Confirm that the watcher channel is closed.
		if _, ok := <-ch; ok {
			t.Fatalf("watcher channel is open on a stopped watcher")
		}
	})
}
