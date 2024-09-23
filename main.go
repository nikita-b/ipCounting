package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // profiling
	"os"
	"time"
)

type AlgorithmType int

const (
	Bitmap AlgorithmType = iota
	BitmapRoaring
	ConcurrentBitmap
)

func timer() func() {
	start := time.Now()
	return func() {
		log.Printf("Execution time is %v\n", time.Since(start))
	}
}

var profile = flag.Bool("profile", false, "Enable profiling")
var filename = flag.String("filename", "", "Path to file with IPs")
var algo = flag.Int("algo", int(Bitmap), "Algorithm to use: 0: Bitmap, 1: BitmapRoaring, 2: ConcurrentBitmap")
var concurrency = flag.Int("concurrency", 5, "Number of workers for concurrent algorithm")

func main() {
	defer timer()()
	flag.Parse()

	if *filename == "" {
		log.Fatalf("Please, provide path to file with IP addresses\n")
		return
	}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}
	defer file.Close()

	if *profile {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalf("Can't get file info: %v\n", err)
	}

	totalFileSize := fileInfo.Size()
	var ipCounter IPCounter
	progress := NewProgressTracker(totalFileSize)
	progress.Start()

	switch *algo {
	case int(Bitmap):
		ipCounter = NewBitmapIPCounter()
	case int(BitmapRoaring):
		ipCounter = NewBitmapRoaringIPCounter()
	case int(ConcurrentBitmap):
		ipCounter = NewBitmapConcurrent(*concurrency)
	default:
		log.Fatalf("Unknown algorithm: %d\n", *algo)
	}

	if *algo != int(ConcurrentBitmap) {
		err = ProcessFile(ipCounter, file, progress)
		if err != nil {
			log.Fatalf("Error processing file: %v\n", err)
		}
	} else {
		err = ProcessFileConcurrency(ipCounter, file, progress, *concurrency)
		if err != nil {
			log.Fatalf("Error processing file: %v\n", err)
		}
	}

	progress.Stop()

	uniqueCount := ipCounter.Count()

	fmt.Printf("Unique IP addresses: %d\n", uniqueCount)
}
