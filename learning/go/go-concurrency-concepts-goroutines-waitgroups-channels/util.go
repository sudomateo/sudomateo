package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// readFn is the type of function used to read all the files within tmpDir,
// sleeping for readDelay between each file read.
type readFn func(tmpDir string, readDelay time.Duration)

// measureTime runs the given fn and prints out how long that function took to
// execute along with other information.
func measureTime(name string, tmpDir string, readDelay time.Duration, fn readFn) {
	fmt.Printf("%s: reading files...\n", name)
	start := time.Now()
	fn(tmpDir, readDelay)
	fmt.Printf("%s: done: took %v\n", name, time.Since(start))
	fmt.Println("---")
}

// setup creates numFiles files populated with some content in a temporary
// directory and returns the name of the temporary directory to the caller.
func setup(numFiles int) (string, error) {
	tmpDir, err := os.MkdirTemp(".", "tmp")
	if err != nil {
		return "", err
	}

	for i := range numFiles {
		if err := os.WriteFile(
			filepath.Join(tmpDir, fmt.Sprintf("file%02d", i)),
			[]byte(fmt.Sprintf("Content for file%02d", i)),
			0644,
		); err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}
