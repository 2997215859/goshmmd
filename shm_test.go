package md

import (
	"encoding/json"
	"fmt"
	"gitlab-dev.qxinvest.com/gomd/md/datatype"
	"gitlab-dev.qxinvest.com/gomd/md/shm"
	"gitlab-dev.qxinvest.com/gomd/md/shmconsumer"
	"gitlab-dev.qxinvest.com/gomd/md/timescale/t092500"
	"testing"
	"time"
	"unsafe"
)

var snapshotCnt uint64 = 0
var orderCnt uint64 = 0
var transCnt uint64 = 0

func CallbackTest(p unsafe.Pointer, dataType uint64) {
	switch dataType {
	case datatype.TypeSnapshot:
		md := datatype.CopySnapshot(shm.GetSnapshot(p))
		b, _ := json.Marshal(md)
		snapshotCnt++
		if snapshotCnt%10000 == 0 {
			fmt.Printf("snapshot: %s\n", string(b))
		}
	case datatype.TypeOrder:
		order := datatype.CopyOrder(shm.GetOrder(p))
		fmt.Printf("order: %+v\n", order)
	case datatype.TypeTransaction:
		transaction := datatype.CopyTransaction(shm.GetTransaction(p))
		fmt.Printf("transaction: %+v\n", transaction)
	case datatype.TypeOrderTransaction:
		bufferType := shm.GetBufferType(p)
		switch bufferType {
		case datatype.TypeOrder:
			order := datatype.CopyOrder(shm.GetUnionOrder(p))
			fmt.Printf("order_transaction.order: %+v, extraType: %+v, order: %+v\n", p, bufferType, order)
		case datatype.TypeTransaction:
			transaction := datatype.CopyTransaction(shm.GetUnionTransaction(p))
			fmt.Printf("order_transaction.transaction: %+v, extraType: %+v, order: %+v\n", p, bufferType, transaction)
		}
	case datatype.TypeOrderTransactionExtra:
		bufferType := shm.GetBufferType(p)
		switch bufferType {
		case datatype.TypeOrderExtra:
			orderExtra := datatype.CopyOrderExtra(shm.GetUnionOrderExtra(p))
			b, _ := json.Marshal(orderExtra)
			orderCnt = orderCnt + 1
			if orderCnt%10000 == 0 {
				fmt.Printf("order_transaction_extra.order_extra: %+v, extraType: %+v, order_extra: %s\n", p, bufferType, string(b))
			}
		case datatype.TypeTransactionExtra:
			transactionExtra := datatype.CopyTransactionExtra(shm.GetUnionTransactionExtra(p))
			b, _ := json.Marshal(transactionExtra)
			transCnt = transCnt + 1
			if transCnt%10000 == 0 {
				fmt.Printf("order_transaction_extra.transaction_extra: %+v, extraType: %+v, transaction_extra: %s\n", p, bufferType, string(b))
			}
		default:
			fmt.Printf("unkonw bufferType: %d\n", bufferType)
		}
	}
}

func SnapshotCallbackTest(snapshot *datatype.Snapshot) {
	b, _ := json.Marshal(snapshot)
	snapshotCnt++
	if snapshotCnt%10000 == 0 {
		fmt.Printf("snapshot: %s\n", string(b))
	}
}

func orderExtraCallbackTest(orderExtra *datatype.OrderExtra) {
	b, _ := json.Marshal(orderExtra)
	orderCnt++
	if orderCnt%10000 == 0 {
		fmt.Printf("orderExtra: %s\n", string(b))
	}
}

func transactionExtraCallbackTest(transactionExtra *datatype.TransactionExtra) {
	b, _ := json.Marshal(transactionExtra)
	transCnt++
	if transCnt%10000 == 0 {
		fmt.Printf("transactionExtra: %s\n", string(b))
	}
}

func tiCallbackTest(ti int) {
	fmt.Printf("ti: %d, time: %s\n", ti, t092500.Ti2Time(ti))
}

func TestAllConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha;/mnt/huge/ha_order_transaction", shmconsumer.WithCallback(CallbackTest))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestOTConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha_order_transaction", shmconsumer.WithCallback(CallbackTest), shmconsumer.WithStart(0))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestSnapshotConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha", shmconsumer.WithSnapshotCallback(SnapshotCallbackTest), shmconsumer.WithStart(0))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestOrderExtraConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha_order_transaction", shmconsumer.WithOrderExtraCallback(orderExtraCallbackTest), shmconsumer.WithStart(1000000))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestTransactionExtraConsumer(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha_order_transaction", shmconsumer.WithTransactionExtraCallback(transactionExtraCallbackTest), shmconsumer.WithStart(1000000))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}

func TestStop(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha", shmconsumer.WithSnapshotCallback(SnapshotCallbackTest), shmconsumer.WithStart(1000000))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	go consumer.Run()
	time.Sleep(10 * time.Second)
	consumer.Stop()
}

func TestTiCallback(t *testing.T) {
	consumer, err := shmconsumer.New("/mnt/huge/ha", shmconsumer.WithTiCallback(tiCallbackTest), shmconsumer.WithStart(0))
	if err != nil {
		t.Errorf("error: %s", err)
	}
	consumer.Run()
}
