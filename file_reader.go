package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func parseIP(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("invalid IPv4 address: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ip), nil
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

func IPAddressReader(reader *bufio.Reader) ([]byte, error) {
	var line []byte
	for {
		chunk, isPrefix, err := reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = append(line, chunk...)
		if !isPrefix {
			break
		}
	}
	return line, nil
}

func ipToUint32(ip []byte) (uint32, error) {
	var result uint32
	var num uint32
	var shift uint = 24

	start := 0
	for i := 0; i <= len(ip); i++ {
		if i == len(ip) || ip[i] == '.' {
			if start == i {
				return 0, fmt.Errorf("empty octet")
			}
			num = 0
			for j := start; j < i; j++ {
				if ip[j] < '0' || ip[j] > '9' {
					return 0, fmt.Errorf("invalid character in IP")
				}
				num = num*10 + uint32(ip[j]-'0')
			}
			if num > 255 {
				return 0, fmt.Errorf("octet value out of range")
			}
			result |= num << shift
			shift -= 8
			start = i + 1
		}
	}
	return result, nil
}

func ProcessFile(ipc IPCounter, file *os.File, progress *ProgressTracker) error {
	reader := bufio.NewReader(file)

	for {
		line, err := IPAddressReader(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		ipInt, err := ipToUint32(line)
		if err != nil {
			fmt.Printf("Error parsing IP %s: %v\n", string(line), err)
			continue
		}
		ipc.Add(&ipInt)
		progress.Increase(int64(len(line)))
	}
	return nil
}

func ProcessFileConcurrency(ipc IPCounter, filePath string, concurrency int) error {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	fileSize := fileStat.Size()
	chunkSize := fileSize / int64(concurrency)

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		startOffset := int64(i) * chunkSize
		var endOffset = startOffset + chunkSize
		wg.Add(1)
		go func(workerID int, startOffset, endOffset int64) {
			defer wg.Done()
			err := processChunk(ipc, filePath, workerID, startOffset, endOffset)
			if err != nil {
				log.Printf("Worker %d has error: %v", workerID, err)
			}
		}(i, startOffset, endOffset)
	}
	wg.Wait()
	return nil
}

func processChunk(ipc IPCounter, filePath string, workerID int, startOffset, endOffset int64) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Seek(startOffset, io.SeekStart)
	if err != nil {
		return err
	}

	reader := bufio.NewReaderSize(file, 2048*2048)
	bytesRead := startOffset

	if startOffset != 0 {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		bytesRead += int64(len(line))
	}

	readerBatchSize := 2048
	batch := make([]uint32, 0, readerBatchSize)

	for {
		if bytesRead >= endOffset {
			break
		}
		line, err := IPAddressReader(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		ipInt, err := ipToUint32(line)
		if err != nil {
			log.Printf("Can't parse IP address: %s. Err: %s", ipInt, err)
			continue
		}
		bytesRead += int64(len(line))

		batch = append(batch, ipInt)
		if len(batch) >= readerBatchSize {
			ipc.AddConcurrent(&batch, workerID)
			batch = make([]uint32, 0, readerBatchSize)
		}

		if err == io.EOF {
			break
		}
	}

	if len(batch) > 0 {
		ipc.AddConcurrent(&batch, workerID)
	}
	return nil
}
