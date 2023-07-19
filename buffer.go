package shmconsumer

import (
	"sync/atomic"
	"unsafe"
)

/*
* Buffer 头信息
 */

const kHugeFileNameSize int32 = 128

type MDGatewayInfo struct {
	totalSize  uint64
	length     int32
	bufferInfo [0]BufferInfo
}

type BufferInfo struct {
	hugeFile  [kHugeFileNameSize]byte
	totalSize uint64
	eachSize  uint64
	length    uint64
	dataType  uint64
}

type Buffer struct {
	totalSize uint64
	mask      uint64
	eachSize  uint64
	dataType  uint64
	padding1  [4]uint64
	tailIndex uint64
	padding2  [7]uint64
	buffer    [0]byte
}

func (buffer *Buffer) GetNextAddrNoBlock(index *uint64) unsafe.Pointer {
	endIndex := atomic.LoadUint64(&buffer.tailIndex)
	if *index == endIndex {
		return nil
	}

	if *index > endIndex {
		*index = endIndex
		return nil
	}
	b := uintptr(unsafe.Pointer(&buffer.buffer)) + uintptr((*index&buffer.mask)*buffer.eachSize)
	p := unsafe.Pointer(&b)
	*index = *index + 1
	//md := *(**MarketData)(unsafe.Pointer(&b))
	return p
}
