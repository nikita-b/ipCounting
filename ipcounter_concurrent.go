package main

import (
	"github.com/RoaringBitmap/roaring/v2"
)

type IPCounterConcurrent struct {
	bitmap      []*roaring.Bitmap
	concurrency int
}

func NewBitmapConcurrent(concurrency int) *IPCounterConcurrent {
	bitmaps := make([]*roaring.Bitmap, concurrency)
	for i := 0; i < concurrency; i++ {
		bitmaps[i] = roaring.NewBitmap()
	}
	return &IPCounterConcurrent{
		bitmap:      bitmaps,
		concurrency: concurrency,
	}
}

func (ipc *IPCounterConcurrent) AddConcurrent(ipAddrRepresentation *[]uint32, workerId int) {
	ipc.bitmap[workerId].AddMany(*ipAddrRepresentation)
}

func (ipc *IPCounterConcurrent) Add(ipAddrRepresentation uint32) {
	// empty
}

func (ipc *IPCounterConcurrent) Count() uint64 {
	combineBitMap := roaring.ParOr(ipc.concurrency, ipc.bitmap...)
	return combineBitMap.GetCardinality()
}
