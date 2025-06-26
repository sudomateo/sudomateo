package main

import (
	"context"
	"runtime"
	"time"
)

// GOMAXPROCSWatcher watches for GOMAXPROCS changes.
type GOMAXPROCSWatcher struct {
	ch   chan int
	done chan struct{}
	cpus int
}

// NewGOMAXPROCSWatcher creates a new GOMAXPROCSWatcher.
func NewGOMAXPROCSWatcher() *GOMAXPROCSWatcher {
	return &GOMAXPROCSWatcher{
		ch:   make(chan int),
		done: make(chan struct{}),
		cpus: runtime.GOMAXPROCS(0),
	}
}

// Watch starts watching for GOMAXPROCS changes.
func (w *GOMAXPROCSWatcher) Watch(ctx context.Context, interval time.Duration) <-chan int {
	go w.watch(ctx, interval)
	return w.ch
}

// watch checks for GOMAXPROCS changes every interval.
func (w *GOMAXPROCSWatcher) watch(ctx context.Context, interval time.Duration) {
	defer close(w.ch)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		case <-ticker.C:
			cpus := runtime.GOMAXPROCS(0)

			// No change. Skip.
			if cpus == w.cpus {
				continue
			}

			w.cpus = cpus
			w.ch <- w.cpus
		}
	}
}

// Stop stops the watcher.
func (w *GOMAXPROCSWatcher) Stop() {
	close(w.done)
}
