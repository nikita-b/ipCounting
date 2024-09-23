package main

import (
	"math/bits"
)

const (
	totalPossibleIpAddresses = 4294967296
	bitesIn64Int             = 64
)

type IPCounter interface {
	Count() uint64
	Add(ipAddr uint32)
	AddConcurrent(ipAddr []uint32, workerId int)
}

type IPCounterBitMap struct {
	bitmap []uint64
}

func NewBitmapIPCounter() *IPCounterBitMap {
	return &IPCounterBitMap{
		bitmap: make([]uint64, totalPossibleIpAddresses/bitesIn64Int),
	}
}

func (ipc *IPCounterBitMap) Add(ipAddrRepresentation uint32) {
	chunkNumber := ipAddrRepresentation / bitesIn64Int
	positionInChunk := ipAddrRepresentation % bitesIn64Int
	ipc.bitmap[chunkNumber] |= 1 << positionInChunk
}

func (ipc *IPCounterBitMap) AddConcurrent(ipAddrRepresentation []uint32, workerId int) {
	// empty
}

func (ipc *IPCounterBitMap) Count() uint64 {
	uniqueCount := 0
	for _, num := range ipc.bitmap {
		uniqueCount += bits.OnesCount64(num)
	}
	return uint64(uniqueCount)
}
