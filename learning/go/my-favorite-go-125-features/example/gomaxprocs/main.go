package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func main() {
	watcher := NewGOMAXPROCSWatcher()
	ch := watcher.Watch(context.Background(), 1*time.Second)
	cpus := runtime.GOMAXPROCS(0)

	for {
		select {
		case updatedCPUs := <-ch:
			fmt.Printf("GOMAXPROCS updated: %02d -> %02d\n", cpus, updatedCPUs)
			cpus = updatedCPUs
		default:
			fmt.Printf("Using %02d CPUs...\n", cpus)
			time.Sleep(3 * time.Second)
		}
	}
}
