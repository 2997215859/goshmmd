package shmconsumer

import (
	"bytes"
	"fmt"
	"gitlab-dev.qxinvest.com/gomd/md/datatype"
	"gitlab-dev.qxinvest.com/gomd/md/shm"
	"gitlab-dev.qxinvest.com/gomd/md/timescale"
	"strings"
	"sync"
	"unsafe"
)

const DefaultBufferSize = 10000000

type Pointer struct {
	pointer  unsafe.Pointer
	dataType uint64
}

type Consumer struct {
	// config
	filePaths       []string
	startIndex      uint64
	startIndexGroup []uint64 // 用于处理多个 buffer 的情况
	bufferSize      uint64

	// data 一个 filepath 可能对应多个 bufferinfo; 每个 bufferInfo 对应一个 buffer
	bufferInfoList []*shm.BufferInfo
	bufferList     []*shm.Buffer
	cursors        []uint64

	stopChan chan struct{}

	callback                 Callback
	snapshotCallback         SnapshotCallback
	orderCallback            OrderCallback
	transactionCallback      TransactionCallback
	orderExtraCallback       OrderExtraCallback
	transactionExtraCallback TransactionExtraCallback

	allChannel              chan *Pointer
	snapshotChannel         chan *datatype.Snapshot
	orderChannel            chan *datatype.Order
	transactionChannel      chan *datatype.Transaction
	orderExtraChannel       chan *datatype.OrderExtra
	transactionExtraChannel chan *datatype.TransactionExtra

	// 一些辅助功能
	latestTimestampMtx sync.Mutex
	latestTimestamp    int64

	latestTiMtx sync.Mutex
	latestTi    int

	tiCallback TiCallback

	tiSeqChannel  chan int
	tiSeqCallback TiCallback

	timescale *timescale.TimeScale
}

func New(filepath string, opts ...Option) (*Consumer, error) {
	filePaths := strings.Split(filepath, ";")
	consumer := &Consumer{
		filePaths:       filePaths,
		startIndex:      DefaultConsumerStartIndex,
		bufferSize:      DefaultBufferSize,
		stopChan:        make(chan struct{}),
		latestTimestamp: -1,
		latestTi:        0,
		timescale:       timescale.DefaultTimeScale,
	}
	for _, o := range opts {
		o(consumer)
	}

	for _, filepath := range filePaths {
		// 1. 映射 buffer 头信息
		bufferHeaderSize := int(unsafe.Sizeof(shm.MDGatewayInfo{}))
		b, err := shm.Alloc(filepath, bufferHeaderSize)
		if err != nil {
			return nil, err
		}

		gatewayInfo := *(**shm.MDGatewayInfo)(unsafe.Pointer(&b))
		if gatewayInfo == nil {
			return nil, fmt.Errorf("error: filepath(%s) gatewayinfo is nil", filepath)
		}

		bufferNum := gatewayInfo.Length
		bufferInfoSize := unsafe.Sizeof(shm.BufferInfo{})
		p := uintptr(unsafe.Pointer(&gatewayInfo.BufferInfo))
		for i := int32(0); i < bufferNum; i++ {
			bufferInfo := *(**shm.BufferInfo)(unsafe.Pointer(&p))
			consumer.bufferInfoList = append(consumer.bufferInfoList, bufferInfo)
			p = p + bufferInfoSize
		}

		for idx, bufferInfo := range consumer.bufferInfoList {
			// 2. 根据 bufferInfo 映射 buffer
			bufferPath := string(bytes.Trim(bufferInfo.HugeFile[:], "\x00"))
			if bufferPath == "" {
				return nil, fmt.Errorf("error: filepath(%s) buffer path is empty", filepath)
			}

			b, err = shm.Alloc(bufferPath, int(unsafe.Sizeof(shm.Buffer{}))+int(bufferInfo.TotalSize))
			if err != nil {
				return nil, err
			}

			buffer := *(**shm.Buffer)(unsafe.Pointer(&b))
			consumer.bufferList = append(consumer.bufferList, buffer)

			// 3. 处理 start index
			startIndex := consumer.startIndex
			if idx < len(consumer.startIndexGroup) {
				startIndex = consumer.startIndexGroup[idx]
			}
			if startIndex == DefaultConsumerStartIndex || startIndex > buffer.TailIndex {
				consumer.cursors = append(consumer.cursors, buffer.TailIndex)
			} else {
				consumer.cursors = append(consumer.cursors, consumer.startIndex)
			}
		}
	}

	return consumer, nil
}

func (c *Consumer) MDCallback() {
	if c.callback != nil {
		c.allChannel = make(chan *Pointer, c.bufferSize)
		go func() {
			for {
				select {
				case data := <-c.allChannel:
					c.callback(data.pointer, data.dataType)
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.snapshotCallback != nil {
		c.snapshotChannel = make(chan *datatype.Snapshot, c.bufferSize)
		go func() {
			for {
				select {
				case snapshot := <-c.snapshotChannel:
					c.snapshotCallback(snapshot)
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.orderCallback != nil {
		c.orderChannel = make(chan *datatype.Order, c.bufferSize)
		go func() {
			for {
				select {
				case order := <-c.orderChannel:
					c.orderCallback(order)
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.transactionCallback != nil {
		c.transactionChannel = make(chan *datatype.Transaction, c.bufferSize)
		go func() {
			for {
				select {
				case transaction := <-c.transactionChannel:
					c.transactionCallback(transaction)
				case <-c.stopChan:
					return
				}
			}
		}()
	}
	if c.orderExtraCallback != nil {
		c.orderExtraChannel = make(chan *datatype.OrderExtra, c.bufferSize)
		go func() {
			for {
				select {
				case order := <-c.orderExtraChannel:
					c.orderExtraCallback(order)
				case <-c.stopChan:
					return
				}
			}
		}()
	}

	if c.transactionExtraCallback != nil {
		c.transactionExtraChannel = make(chan *datatype.TransactionExtra, c.bufferSize)
		go func() {
			for {
				select {
				case transactionExtra := <-c.transactionExtraChannel:
					c.transactionExtraCallback(transactionExtra)
				case <-c.stopChan:
					return
				}
			}
		}()
	}

	if c.tiSeqCallback != nil {
		c.tiSeqChannel = make(chan int, 1024)
		go func() {
			for {
				select {
				case ti := <-c.tiSeqChannel:
					c.tiSeqCallback(ti)
				case <-c.stopChan:
					return
				}
			}
		}()
	}
}

func (c *Consumer) ReadBuffer(buffer *shm.Buffer, currentIndex *uint64) {
	datapoint := buffer.GetNextAddrNoBlock(currentIndex)
	if datapoint == nil {
		return
	}

	if c.callback != nil && c.allChannel != nil {
		c.allChannel <- &Pointer{
			pointer:  datapoint,
			dataType: buffer.DataType,
		}
	}

	switch buffer.DataType {
	case datatype.TypeSnapshot:
		snapshot := datatype.CopySnapshot(shm.GetSnapshot(datapoint))
		if snapshot != nil {
			go c.SetTimestamp(snapshot.TimestampS)
			c.CallTimer(snapshot.UpdateTime)
		}
		if snapshot != nil && c.snapshotCallback != nil && c.snapshotChannel != nil {
			c.snapshotChannel <- snapshot
		}
	case datatype.TypeOrder:
		order := datatype.CopyOrder(shm.GetOrder(datapoint))
		if order != nil && c.orderCallback != nil && c.orderChannel != nil {
			c.orderChannel <- order
		}
	case datatype.TypeTransaction:
		transaction := datatype.CopyTransaction(shm.GetTransaction(datapoint))
		if transaction != nil && c.transactionCallback != nil && c.transactionChannel != nil {
			c.transactionChannel <- transaction
		}
	case datatype.TypeOrderTransaction:
		bufferType := shm.GetBufferType(datapoint)
		switch bufferType {
		case datatype.TypeOrder:
			order := datatype.CopyOrder(shm.GetUnionOrder(datapoint))
			if order != nil && c.orderCallback != nil && c.orderChannel != nil {
				c.orderChannel <- order
			}
		case datatype.TypeTransaction:
			transaction := datatype.CopyTransaction(shm.GetUnionTransaction(datapoint))
			if transaction != nil && c.transactionCallback != nil && c.transactionChannel != nil {
				c.transactionChannel <- transaction
			}
		}
	case datatype.TypeOrderTransactionExtra:
		bufferType := shm.GetBufferType(datapoint)
		switch bufferType {
		case datatype.TypeOrderExtra:
			orderExtra := datatype.CopyOrderExtra(shm.GetUnionOrderExtra(datapoint))
			if orderExtra != nil && c.orderExtraCallback != nil && c.orderExtraChannel != nil {
				c.orderExtraChannel <- orderExtra
			}
		case datatype.TypeTransactionExtra:
			transactionExtra := datatype.CopyTransactionExtra(shm.GetUnionTransactionExtra(datapoint))
			if transactionExtra != nil && c.transactionExtraCallback != nil && c.transactionExtraChannel != nil {
				c.transactionExtraChannel <- transactionExtra
			}
		}
	}
}

func (c *Consumer) LoopBuffer(idx int, buffer *shm.Buffer) {
	for {
		select {
		case <-c.stopChan:
			return
		default:
			c.ReadBuffer(buffer, &c.cursors[idx])
		}
	}
}

func (c *Consumer) Run() {

	c.MDCallback()

	var wg sync.WaitGroup
	wg.Add(len(c.bufferList))

	for idx, buffer := range c.bufferList {
		go func(idx int, buffer *shm.Buffer) {
			c.LoopBuffer(idx, buffer)
			wg.Done()
		}(idx, buffer)
	}

	wg.Wait()
}

func (c *Consumer) Start() {
	go c.Run()
}

func (c *Consumer) Stop() {
	close(c.stopChan)
}

// 一些辅助性的功能

func (c *Consumer) SetTimestamp(timestamp int64) {
	c.latestTimestampMtx.Lock()
	defer c.latestTimestampMtx.Unlock()

	c.latestTimestamp = timestamp
}

func (c *Consumer) GetTimestamp() int64 {
	c.latestTimestampMtx.Lock()
	defer c.latestTimestampMtx.Unlock()

	return c.latestTimestamp
}

func (c *Consumer) UpdateLatestTi(ti int) bool {
	c.latestTiMtx.Lock()
	defer c.latestTiMtx.Unlock()

	if ti > c.latestTi {
		c.latestTi = ti
		return true
	}
	return false
}

func (c *Consumer) CallTimer(updateTime string) {
	ti := c.timescale.GetTi(updateTime)

	if ti < 0 {
		return
	}

	if c.UpdateLatestTi(ti) { // 如果更新成功，则调用 1 分钟回调
		if c.tiCallback != nil {
			go c.tiCallback(ti)
		}
		if c.tiSeqCallback != nil {
			c.tiSeqChannel <- ti
		}
	}
}
