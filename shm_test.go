package shmconsumer

import (
	"encoding/json"
	"fmt"
	"gitlab-dev.qxinvest.com/gomd/shmconsumer/raw"
	"testing"
	"unsafe"
)

func MDCallback(p unsafe.Pointer, dataType uint64) {
	switch dataType {
	case raw.TypeMarketData:
		md := CopyMarketData(raw.GetMarketData(p))
		b, _ := json.Marshal(md)
		fmt.Printf("market_data: %s\n", string(b))
	case raw.TypeOrder:
		order := CopyOrder(raw.GetOrder(p))
		fmt.Printf("order: %+v\n", order)
	case raw.TypeTransaction:
		transaction := CopyTransaction(raw.GetTransaction(p))
		fmt.Printf("order: %+v\n", transaction)
	case raw.TypeOrderTransaction:
		bufferType := raw.GetBufferType(p)
		switch bufferType {
		case raw.TypeOrder:
			order := CopyOrder(raw.GetUnionOrder(p))
			fmt.Printf("order_transaction.order: %+v, extraType: %+v, order: %+v\n", p, bufferType, order)
		case raw.TypeTransaction:
			transaction := CopyTransaction(raw.GetUnionTransaction(p))
			fmt.Printf("order_transaction.transaction: %+v, extraType: %+v, order: %+v\n", p, bufferType, transaction)
		}
	case raw.TypeOrderTransactionExtra:
		bufferType := raw.GetBufferType(p)
		switch bufferType {
		case raw.TypeOrderExtra:
			orderExtra := CopyOrderExtra(raw.GetUnionOrderExtra(p))
			b, _ := json.Marshal(orderExtra)
			fmt.Printf("order_transaction_extra.order_extra: %+v, extraType: %+v, order_extra: %s\n", p, bufferType, string(b))
		case raw.TypeTransaction:
			transactionExtra := CopyTransactionExtra(raw.GetUnionTransactionExtra(p))
			b, _ := json.Marshal(transactionExtra)
			fmt.Printf("order_transaction_extra.transaction_extra: %+v, extraType: %+v, order: %s\n", p, bufferType, string(b))
		}
	}
}

func TestMDConsumer(t *testing.T) {
	consumer, err := New("/mnt/huge/ha", WithCallback(MDCallback), WithStart(10000000))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestOTConsumer(t *testing.T) {
	consumer, err := New("/mnt/huge/ha_order_transaction", WithCallback(MDCallback), WithStart(10000000))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}
