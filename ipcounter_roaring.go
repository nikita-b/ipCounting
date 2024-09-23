package main

import (
	"github.com/RoaringBitmap/roaring/v2"
)

type IPCounterRoaring struct {
	bitmap *roaring.Bitmap
}

func NewBitmapRoaringIPCounter() *IPCounterRoaring {
	return &IPCounterRoaring{
		bitmap: roaring.NewBitmap(),
	}
}

func (ipc *IPCounterRoaring) Add(ipAddrRepresentation uint32) {
	ipc.bitmap.Add(ipAddrRepresentation)
}

func (ipc *IPCounterRoaring) Count() uint64 {
	return ipc.bitmap.GetCardinality()
}
