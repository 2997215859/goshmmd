package shmconsumer

import (
	"math"
	"unsafe"
)

const DefaultConsumerStartIndex uint64 = math.MaxUint64

type Callback func(data unsafe.Pointer, dataType uint64)

type Option func(consumer *Consumer)

func WithStart(startIndex uint64) Option {
	return func(consumer *Consumer) {
		consumer.startIndex = startIndex
	}
}

func WithCallback(callback Callback) Option {
	return func(consumer *Consumer) {
		consumer.callback = callback
	}
}
