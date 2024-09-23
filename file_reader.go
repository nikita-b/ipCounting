package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func parseIP(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IPv4: %s", ipStr)
	}
	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ipv4), nil
}

//func parseIP(ip string) (uint32, error) {
//	ipParts := strings.Split(ip, ".")
//	if len(ipParts) != 4 {
//		log.Printf("String is not IPv4 address")
//		return 0, errors.New("invalid IP")
//	}
//	var ipInt uint64
//	for i, part := range ipParts {
//		octet, err := strconv.Atoi(part)
//		if err != nil {
//			return 0, fmt.Errorf("invalid octet %d: %v", i+1, err)
//		}
//		if octet < 0 || octet > 255 {
//			return 0, fmt.Errorf("octet %d out of range: %d", i+1, octet)
//		}
//		ipInt = ipInt*256 + uint64(octet)
//	}
//	return ipInt, nil
//}

func ProcessFile(ipc IPCounter, file *os.File, progress *ProgressTracker) error {
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			return nil
		}
		progress.Increase(int64(len(line)))

		ipAddrStr := strings.Trim(line, "\n")
		ipAddr, err := parseIP(ipAddrStr)
		if err != nil {
			log.Printf("Can't parse IP address: %s. Err: %s", ipAddrStr, err)
			continue
		}
		ipc.Add(ipAddr)
	}
}

func ProcessFileConcurrency(ipc IPCounter, file *os.File, progress *ProgressTracker, concurrency int) error {
	ipAddrQueue := make(chan uint32, 10000) // ~4MB

	go func() {
		reader := bufio.NewReaderSize(file, 1<<10)
		defer close(ipAddrQueue)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				log.Fatalf("Error reading file: %v", err)
			}
			if err == io.EOF {
				break
			}
			progress.bytesRead = progress.bytesRead + int64(len(line))

			ipAddrStr := strings.Trim(line, "\n")
			ipAddr, err := parseIP(ipAddrStr)
			if err != nil {
				log.Printf("Can't parse IP address: %s. Err: %s", ipAddrStr, err)
				continue
			}
			ipAddrQueue <- ipAddr
		}
	}()
	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := make([]uint32, 0, 100) // Batch size of 100
			for ipInt := range ipAddrQueue {
				batch = append(batch, ipInt)
				if len(batch) >= 100 {
					ipc.AddConcurrent(batch, i)
					batch = make([]uint32, 0, 100)
				}
			}
			if len(batch) > 0 {
				ipc.AddConcurrent(batch, i)
			}
		}()
	}
	wg.Wait()
	return nil
}
