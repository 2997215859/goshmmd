package shm

import (
	"sync/atomic"
	"unsafe"
)

/*
* Buffer 头信息
 */

const kHugeFileNameSize int32 = 128

type MDGatewayInfo struct {
	TotalSize  uint64
	Length     int32
	BufferInfo [0]BufferInfo
}

type BufferInfo struct {
	HugeFile  [kHugeFileNameSize]byte
	TotalSize uint64
	EachSize  uint64
	Length    uint64
	DataType  uint64
}

type Buffer struct {
	TotalSize uint64
	Mask      uint64
	EachSize  uint64
	DataType  uint64
	Padding1  [4]uint64
	TailIndex uint64
	Padding2  [7]uint64
	Buffer    [0]byte
}

func (buffer *Buffer) GetNextAddrNoBlock(index *uint64) unsafe.Pointer {
	endIndex := atomic.LoadUint64(&buffer.TailIndex)
	if *index == endIndex {
		return nil
	}

	if *index > endIndex {
		*index = endIndex
		return nil
	}
	b := uintptr(unsafe.Pointer(&buffer.Buffer)) + uintptr((*index&buffer.Mask)*buffer.EachSize)
	p := unsafe.Pointer(&b)
	*index = *index + 1
	//md := *(**MarketData)(unsafe.Pointer(&b))
	return p
}
