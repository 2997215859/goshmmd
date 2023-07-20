package md

import (
	"encoding/json"
	"fmt"
	"gitlab-dev.qxinvest.com/gomd/md/shm"
	"gitlab-dev.qxinvest.com/gomd/md/shmconsumer"
	"testing"
	"unsafe"
)

var orderCnt uint64 = 0
var transCnt uint64 = 0

func CallbackTest(p unsafe.Pointer, dataType uint64) {
	switch dataType {
	case TypeMarketData:
		md := CopyMarketData(shm.GetMarketData(p))
		b, _ := json.Marshal(md)
		fmt.Printf("market_data: %s\n", string(b))
	case TypeOrder:
		order := CopyOrder(shm.GetOrder(p))
		fmt.Printf("order: %+v\n", order)
	case TypeTransaction:
		transaction := CopyTransaction(shm.GetTransaction(p))
		fmt.Printf("transaction: %+v\n", transaction)
	case TypeOrderTransaction:
		bufferType := shm.GetBufferType(p)
		switch bufferType {
		case TypeOrder:
			order := CopyOrder(shm.GetUnionOrder(p))
			fmt.Printf("order_transaction.order: %+v, extraType: %+v, order: %+v\n", p, bufferType, order)
		case TypeTransaction:
			transaction := CopyTransaction(shm.GetUnionTransaction(p))
			fmt.Printf("order_transaction.transaction: %+v, extraType: %+v, order: %+v\n", p, bufferType, transaction)
		}
	case TypeOrderTransactionExtra:
		bufferType := shm.GetBufferType(p)
		switch bufferType {
		case TypeOrderExtra:
			orderExtra := CopyOrderExtra(shm.GetUnionOrderExtra(p))
			b, _ := json.Marshal(orderExtra)
			orderCnt = orderCnt + 1
			if orderCnt%10000 == 0 {
				fmt.Printf("order_transaction_extra.order_extra: %+v, extraType: %+v, order_extra: %s\n", p, bufferType, string(b))
			}
		case TypeTransactionExtra:
			transactionExtra := CopyTransactionExtra(shm.GetUnionTransactionExtra(p))
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
	consumer, err := shmconsumer.New("/mnt/huge/ha", shmconsumer.WithCallback(CallbackTest))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestOTConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha_order_transaction", shmconsumer.WithCallback(CallbackTest))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}
