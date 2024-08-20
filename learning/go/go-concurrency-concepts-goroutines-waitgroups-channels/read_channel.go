package main

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// readChannel reads all of the files within tmpDir concurrently using a
// channel with no upper bound on its concurrency. Before each file is read
// this function will sleep for readDelay. The content of each file is printed
// to standard output.
func readChannel(tmpDir string, readDelay time.Duration) {
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	doneCh := make(chan struct{})

	for _, entry := range entries {
		go func(entry fs.DirEntry) {
			// Ensure this goroutine tells the receiving channel it's finished
			// executing.
			defer func() {
				doneCh <- struct{}{}
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
	for range entries {
		<-doneCh
	}
}
