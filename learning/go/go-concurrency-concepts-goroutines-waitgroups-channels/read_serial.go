package main

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// readSerial reads all of the files within tmpDir serially. Before each file
// is read this function will sleep for readDelay. The content of each file is
// printed to standard output.
func readSerial(tmpDir string, readDelay time.Duration) {
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		time.Sleep(readDelay)

		b, err := os.ReadFile(filepath.Join(tmpDir, entry.Name()))
		if err != nil {
			panic(err)
		}

		io.WriteString(os.Stdout, string(b)+"\n")
	}
}
