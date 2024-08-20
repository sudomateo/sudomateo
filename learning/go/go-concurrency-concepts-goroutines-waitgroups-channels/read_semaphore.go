package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// readSemaphore reads all of the files within tmpDir concurrently using a wait
// group and semaphore channel. The semaphore channel enforces an upper bound
// on the number of goroutines allowed to run concurrently. Before each file is
// read this function will sleep for readDelay. The content of each file is
// printed to standard output.
func readSemaphore(tmpDir string, readDelay time.Duration) {
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	// Only allow a fixed number of goroutines to run concurrently.
	semaphore := make(chan struct{}, 20)

	for _, entry := range entries {
		// Wait until one of the concurrency slots opens up.
		semaphore <- struct{}{}

		// Tell the wait group we're about to do 1 unit of work.
		wg.Add(1)

		go func(entry fs.DirEntry) {
			// Ensure this goroutine frees up a concurrency slot and tells the
			// wait group the 1 unit of work is finished.
			defer func() {
				<-semaphore
				wg.Done()
			}()

			time.Sleep(readDelay)

			b, err := os.ReadFile(filepath.Join(tmpDir, entry.Name()))
			if err != nil {
				panic(err)
			}

			io.WriteString(os.Stdout, string(b)+"\n")
		}(entry)
	}

	// Wait for all goroutines to finish.
	wg.Wait()
}
