package main

import (
	"github.com/RoaringBitmap/roaring/v2/roaring64"
)

type IPCounterConcurrent struct {
	bitmap      []*roaring64.Bitmap
	concurrency int
}

func NewBitmapConcurrent(concurrency int) *IPCounterConcurrent {
	bitmaps := make([]*roaring64.Bitmap, concurrency)
	for i := 0; i < concurrency; i++ {
		bitmaps[i] = roaring64.NewBitmap()
	}
	return &IPCounterConcurrent{
		bitmap:      bitmaps,
		concurrency: concurrency,
	}
}

func (ipc *IPCounterConcurrent) AddConcurrent(ipAddrRepresentation *[]uint64, workerId int) {
	ipc.bitmap[workerId].AddMany(*ipAddrRepresentation)
}

func (ipc *IPCounterConcurrent) Add(ipAddrRepresentation uint64) {
	// empty
}

func (ipc *IPCounterConcurrent) Count() uint64 {
	combineBitMap := roaring64.ParOr(ipc.concurrency, ipc.bitmap...)
	return combineBitMap.GetCardinality()
}
