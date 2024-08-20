package main

import (
	"os"
	"time"
)

func main() {
	numFiles := 20
	readDelay := 5 * time.Millisecond

	tmpDir, err := setup(numFiles)
	if err != nil {
		panic(err)
	}

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	measureTime("serial", tmpDir, readDelay, readSerial)
	measureTime("waitgroup", tmpDir, readDelay, readWaitGroup)
	measureTime("channel", tmpDir, readDelay, readChannel)
	measureTime("semaphore", tmpDir, readDelay, readSemaphore)
}
