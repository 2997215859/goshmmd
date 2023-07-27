package shmconsumer

import (
	"gitlab-dev.qxinvest.com/gomd/md/datatype"
	"gitlab-dev.qxinvest.com/gomd/md/timescale"
	"math"
	"unsafe"
)

const DefaultConsumerStartIndex uint64 = math.MaxUint64

type Callback func(data unsafe.Pointer, dataType uint64)
type SnapshotCallback func(snapshot *datatype.Snapshot)
type OrderCallback func(order *datatype.Order)
type TransactionCallback func(trans *datatype.Transaction)
type OrderExtraCallback func(orderExtra *datatype.OrderExtra)
type TransactionExtraCallback func(transactionExtra *datatype.TransactionExtra)

type TiCallback func(int)

type Option func(consumer *Consumer)

func WithStart(startIndex uint64) Option {
	return func(consumer *Consumer) {
		consumer.startIndex = startIndex
	}
}

func WithStarts(starts []uint64) Option {
	return func(consumer *Consumer) {
		consumer.startIndexGroup = starts
	}
}

func WithBufferSize(bufferSize uint64) Option {
	return func(consumer *Consumer) {
		consumer.bufferSize = bufferSize
	}
}

func WithCallback(cb Callback) Option {
	return func(consumer *Consumer) {
		consumer.callback = cb
	}
}

func WithSnapshotCallback(cb SnapshotCallback) Option {
	return func(consumer *Consumer) {
		consumer.snapshotCallback = cb
	}
}
func WithOrderCallback(cb OrderCallback) Option {
	return func(consumer *Consumer) {
		consumer.orderCallback = cb
	}
}

func WithTransactionCallback(cb TransactionCallback) Option {
	return func(consumer *Consumer) {
		consumer.transactionCallback = cb
	}
}

func WithOrderExtraCallback(cb OrderExtraCallback) Option {
	return func(consumer *Consumer) {
		consumer.orderExtraCallback = cb
	}
}

func WithTransactionExtraCallback(cb TransactionExtraCallback) Option {
	return func(consumer *Consumer) {
		consumer.transactionExtraCallback = cb
	}
}

func WithTiCallback(cb TiCallback) Option {
	return func(consumer *Consumer) {
		consumer.tiCallback = cb
	}
}

func WithTiSeqCallback(cb TiCallback) Option {
	return func(consumer *Consumer) {
		consumer.tiSeqCallback = cb
	}
}

func WithTimescale(scale timescale.TimeScale) Option {
	return func(consumer *Consumer) {
		consumer.timescale = &scale
	}
}
