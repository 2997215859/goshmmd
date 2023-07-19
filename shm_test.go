package shmconsumer

import (
	"encoding/json"
	"fmt"
	"gitlab-dev.qxinvest.com/gomd/shmconsumer/raw"
	"testing"
	"unsafe"
)

var orderCnt uint64 = 0
var transCnt uint64 = 0

func CallbackTest(p unsafe.Pointer, dataType uint64) {
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
			orderCnt = orderCnt + 1
			if orderCnt%10000 == 0 {
				fmt.Printf("order_transaction_extra.order_extra: %+v, extraType: %+v, order_extra: %s\n", p, bufferType, string(b))
			}
		case raw.TypeTransactionExtra:
			transactionExtra := CopyTransactionExtra(raw.GetUnionTransactionExtra(p))
			b, _ := json.Marshal(transactionExtra)
			transCnt = transCnt + 1
			if transCnt%10000 == 0 {
				fmt.Printf("order_transaction_extra.transaction_extra: %+v, extraType: %+v, order: %s\n", p, bufferType, string(b))
			}
		default:
			fmt.Printf("unkonw bufferType: %d\n", bufferType)
		}
	}
}

func TestMDConsumer(t *testing.T) {
	consumer, err := New("/mnt/huge/ha", WithCallback(CallbackTest), WithStart(0))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestOTConsumer(t *testing.T) {
	consumer, err := New("/mnt/huge/ha_order_transaction", WithCallback(CallbackTest), WithStart(0))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}
