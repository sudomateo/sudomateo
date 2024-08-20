package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// readWaitGroup reads all of the files within tmpDir concurrently using a wait
// group with no upper bound on its concurrency. Before each file is read this
// function will sleep for readDelay. The content of each file is printed to
// standard output.
func readWaitGroup(tmpDir string, readDelay time.Duration) {
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, entry := range entries {
		// Tell the wait group we're about to do 1 unit of work.
		wg.Add(1)

		go func(entry fs.DirEntry) {
			// Ensure this goroutine tells the wait group the 1 unit of work is
			// finished.
			defer wg.Done()

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
