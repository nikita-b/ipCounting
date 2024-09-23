package main

import (
	"math/bits"
	"sync/atomic"
)

type IPCounterConcurrent struct {
	bitmap []uint64
}

func NewBitmapConcurrent() *IPCounterConcurrent {
	return &IPCounterConcurrent{
		bitmap: make([]uint64, totalPossibleIpAddresses/bitesIn64Int),
	}
}

func (ipc *IPCounterConcurrent) Add(ipAddrRepresentation uint32) {
	chunkNumber := ipAddrRepresentation / bitesIn64Int
	positionInChunk := ipAddrRepresentation % bitesIn64Int
	mask := uint64(1) << positionInChunk

	addr := &ipc.bitmap[chunkNumber]
	atomic.AddUint64(addr, mask&^atomic.LoadUint64(addr)&mask)
}

func (ipc *IPCounterConcurrent) Count() uint64 {
	uniqueCount := 0
	for _, num := range ipc.bitmap {
		uniqueCount += bits.OnesCount64(num)
	}
	return uint64(uniqueCount)
}
