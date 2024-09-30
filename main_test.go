package main

import (
	"log"
	"os"
	"strconv"
	"testing"
)

func getenv(key string, defaultValue int) uint64 {
	strValue := os.Getenv(key)
	if len(strValue) == 0 {
		return uint64(defaultValue)
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		panic(err)
	}
	return uint64(value)
}

var expectedFromENV = getenv("GENERATED_IPS", 30000) / 100 * 99

func TestIPCounterBasic(t *testing.T) {
	expectedUniqueCount := uint64(4)
	file, err := os.Open("resources/ip_addresses_basic.txt")
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}
	defer file.Close()

	var ipc = NewBitmapIPCounter()
	var pt = NewProgressTracker(0)
	err = ProcessFile(ipc, file, pt)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}

func TestIPCounterRoaring(t *testing.T) {
	expectedUniqueCount := uint64(4)
	fileName, err := os.Open("resources/ip_addresses_basic.txt")
	var ipc = NewBitmapRoaringIPCounter()
	var pt = NewProgressTracker(0)
	err = ProcessFile(ipc, fileName, pt)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}

func TestIPCounterConcurrentBasic(t *testing.T) {
	expectedUniqueCount := uint64(4)
	var ipc = NewBitmapConcurrent(2)
	var err = ProcessFileConcurrency(ipc, "resources/ip_addresses_basic.txt", 2)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}

func TestIPCounterFull(t *testing.T) {
	expectedUniqueCount := expectedFromENV
	file, err := os.Open("resources/generated_ip_addresses.txt")
	if err != nil {
		log.Fatalf("Can't open file: %v\n", err)
	}
	defer file.Close()

	var ipc = NewBitmapIPCounter()
	var pt = NewProgressTracker(0)
	err = ProcessFile(ipc, file, pt)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}

func TestIPCounterRoaringFull(t *testing.T) {
	expectedUniqueCount := expectedFromENV
	fileName, err := os.Open("resources/generated_ip_addresses.txt")
	var ipc = NewBitmapRoaringIPCounter()
	var pt = NewProgressTracker(0)
	err = ProcessFile(ipc, fileName, pt)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}

func TestIPCounterConcurrentFull(t *testing.T) {
	expectedUniqueCount := expectedFromENV
	var ipc = NewBitmapConcurrent(2)
	var err = ProcessFileConcurrency(ipc, "resources/generated_ip_addresses.txt", 2)
	if err != nil {
		t.Fatalf("Error processing file: %v", err)
	}
	uniqueCount := ipc.Count()

	if uniqueCount != expectedUniqueCount {
		t.Errorf("Expected unique count %d, got %d", expectedUniqueCount, uniqueCount)
	}
}
