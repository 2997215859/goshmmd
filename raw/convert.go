package raw

import (
	"unsafe"
)

// 处理正常类型

func GetMarketData(data unsafe.Pointer) *MarketData {
	return *(**MarketData)(data)
}

func GetOrder(data unsafe.Pointer) *Order {
	return *(**Order)(data)
}

func GetTransaction(data unsafe.Pointer) *Transaction {
	return *(**Transaction)(data)
}

// 处理 extra

func GetOrderExtra(data unsafe.Pointer) *OrderExtra {
	return *(**OrderExtra)(data)
}

func GetTransactionExtra(data unsafe.Pointer) *TransactionExtra {
	return *(**TransactionExtra)(data)
}

// 处理 正常 union 类型

func GetBufferType(data unsafe.Pointer) int {
	return **(**int)(data)
}

func GetUnionOrder(data unsafe.Pointer) *Order {
	union := *(**UnionOrder)(data)
	return &union.order
}

func GetUnionTransaction(data unsafe.Pointer) *Transaction {
	union := *(**UnionTransaction)(data)
	return &union.transaction
}

// 处理 extra union 类型

func GetUnionOrderExtra(data unsafe.Pointer) *OrderExtra {
	union := *(**UnionOrderExtra)(data)
	return &union.orderExtra
}

func GetUnionTransactionExtra(data unsafe.Pointer) *TransactionExtra {
	union := *(**UnionTransactionExtra)(data)
	return &union.transactionExtra
}
