package shmconsumer

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type Consumer struct {
	// config
	filePaths       []string
	startIndex      uint64
	startIndexGroup []uint64 // 用于处理多个 buffer 的情况

	// data 一个 filepath 可能对应多个 bufferinfo; 每个 bufferInfo 对应一个 buffer
	bufferInfoList []*BufferInfo
	bufferList     []*Buffer
	cursors        []uint64
	callback       Callback
}

func New(filepath string, opts ...Option) (*Consumer, error) {
	filePaths := strings.Split(filepath, ",")
	consumer := &Consumer{
		filePaths:  filePaths,
		startIndex: DefaultConsumerStartIndex,
	}
	for _, o := range opts {
		o(consumer)
	}

	for _, filepath := range filePaths {
		// 1. 映射 buffer 头信息
		bufferHeaderSize := int(unsafe.Sizeof(MDGatewayInfo{}))
		b, err := Alloc(filepath, bufferHeaderSize)
		if err != nil {
			return nil, err
		}

		gatewayInfo := *(**MDGatewayInfo)(unsafe.Pointer(&b))
		if gatewayInfo == nil {
			return nil, fmt.Errorf("error: filepath(%s) gatewayinfo is nil", filepath)
		}

		bufferNum := gatewayInfo.length
		bufferInfoSize := unsafe.Sizeof(BufferInfo{})
		p := uintptr(unsafe.Pointer(&gatewayInfo.bufferInfo))
		for i := int32(0); i < bufferNum; i++ {
			bufferInfo := *(**BufferInfo)(unsafe.Pointer(&p))
			consumer.bufferInfoList = append(consumer.bufferInfoList, bufferInfo)
			p = p + bufferInfoSize
		}

		for idx, bufferInfo := range consumer.bufferInfoList {
			// 2. 根据 bufferInfo 映射 buffer
			bufferPath := string(bytes.Trim(bufferInfo.hugeFile[:], "\x00"))
			if bufferPath == "" {
				return nil, fmt.Errorf("error: filepath(%s) buffer path is empty", filepath)
			}

			b, err = Alloc(bufferPath, int(unsafe.Sizeof(Buffer{}))+int(bufferInfo.totalSize))
			if err != nil {
				return nil, err
			}

			buffer := *(**Buffer)(unsafe.Pointer(&b))
			consumer.bufferList = append(consumer.bufferList, buffer)

			// 3. 处理 start index
			startIndex := consumer.startIndex
			if idx < len(consumer.startIndexGroup) {
				startIndex = consumer.startIndexGroup[idx]
			}
			if startIndex == DefaultConsumerStartIndex || startIndex > buffer.tailIndex {
				consumer.cursors = append(consumer.cursors, buffer.tailIndex)
			} else {
				consumer.cursors = append(consumer.cursors, consumer.startIndex)
			}
		}
	}

	return consumer, nil
}

func Alloc(filepath string, size int) ([]byte, error) {
	f, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("filepath error: %s", err)
	}
	defer f.Close()

	b, err := syscall.Mmap(int(f.Fd()), 0, size, syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("mmap(%s) size(%v) error: %s", filepath, size, err)
	}
	if b == nil {
		return nil, fmt.Errorf("error: mmap(%s) data is nil", filepath)
	}
	return b, nil
}

func (c *Consumer) Run() {
	for {
		for idx, buffer := range c.bufferList {
			datapoint := buffer.GetNextAddrNoBlock(&c.cursors[idx])
			if datapoint != nil && c.callback != nil {
				c.callback(datapoint, buffer.dataType)
			}
		}
	}
}
