package main

import (
	"github.com/RoaringBitmap/roaring/v2"
)

type IPCounterConcurrent struct {
	bitmap      []roaring.Bitmap
	concurrency int
}

func NewBitmapConcurrent(concurrency int) *IPCounterConcurrent {
	return &IPCounterConcurrent{
		bitmap:      make([]roaring.Bitmap, concurrency),
		concurrency: concurrency,
	}
}

func (ipc *IPCounterConcurrent) AddConcurrent(ipAddrRepresentation []uint32, workerId int) {
	ipc.bitmap[workerId].AddMany(ipAddrRepresentation)
}

func (ipc *IPCounterConcurrent) Add(ipAddrRepresentation uint32) {
	// empty
}

func (ipc *IPCounterConcurrent) Count() uint64 {
	combineBitMap := roaring.ParOr(0, &ipc.bitmap[0], &ipc.bitmap[1], &ipc.bitmap[2], &ipc.bitmap[3], &ipc.bitmap[4])
	return combineBitMap.GetCardinality()
}
