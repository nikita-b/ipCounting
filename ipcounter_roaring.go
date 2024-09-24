package main

import (
	"github.com/RoaringBitmap/roaring/v2/roaring64"
)

type IPCounterRoaring struct {
	bitmap *roaring64.Bitmap
}

func NewBitmapRoaringIPCounter() *IPCounterRoaring {
	return &IPCounterRoaring{
		bitmap: roaring64.NewBitmap(),
	}
}

func (ipc *IPCounterRoaring) Add(ipAddrRepresentation uint64) {
	ipc.bitmap.Add(ipAddrRepresentation)
}

func (ipc *IPCounterRoaring) AddConcurrent(ipAddrRepresentation *[]uint64, workerId int) {
	// empty
}

func (ipc *IPCounterRoaring) Count() uint64 {
	return ipc.bitmap.GetCardinality()
}
