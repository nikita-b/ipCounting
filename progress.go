package main

import (
	"log"
	"time"
)

type ProgressTracker struct {
	totalSize           int64
	bytesRead           int64
	lastPrintedProgress int64
	done                chan struct{}
}

func NewProgressTracker(totalSize int64) *ProgressTracker {
	return &ProgressTracker{
		totalSize: totalSize,
		done:      make(chan struct{}),
	}
}

func (pt *ProgressTracker) Start() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-pt.done:
				return
			case <-ticker.C:
				progress := pt.bytesRead * 100 / pt.totalSize
				if progress != pt.lastPrintedProgress {
					log.Printf("\rProcessing: %d%% complete\n", progress)
					pt.lastPrintedProgress = progress
				}
			}
		}
	}()
}

func (pt *ProgressTracker) Increase(bytesRead int64) {
	pt.bytesRead = pt.bytesRead + bytesRead
}

func (pt *ProgressTracker) Stop() {
	close(pt.done)
	log.Printf("\rProcessing complete\n")
}
