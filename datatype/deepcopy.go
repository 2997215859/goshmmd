package datatype

import (
	"gitlab-dev.qxinvest.com/gomd/md/shm"
	"math"
	"strings"
)

// ClearString
// []byte 中间含有特殊字符'\0', golang 无法直接识别处理，需要手动截断
func ClearString(s []byte) string {
	idx := strings.Index(string(s), string([]byte{0}))
	if idx == -1 {
		return string(s)
	}
	return string(s[:idx])
}

func Suffix(instrumentId string) string {
	idx := strings.Index(instrumentId, ".")
	if idx == -1 {
		return ""
	}
	return instrumentId[idx+1:]
}

func RemoveNan(arr []float64) {
	for idx, price := range arr {
		if math.IsNaN(price) {
			arr[idx] = 0
		}
	}
}

// md 的自定义序列化

func CopySnapshot(src *shm.Snapshot) *Snapshot {
	askPrice := src.AskPrice
	RemoveNan(askPrice[:])

	askVolume := src.AskVolume
	RemoveNan(askVolume[:])

	bidPrice := src.BidPrice
	RemoveNan(bidPrice[:])

	bidVolume := src.BidVolume
	RemoveNan(bidVolume[:])

	averageAskPrice := src.AverageAskPrice
	if math.IsNaN(averageAskPrice) {
		averageAskPrice = 0
	}

	averageBidPrice := src.AverageBidPrice
	if math.IsNaN(averageBidPrice) {
		averageBidPrice = 0
	}

	totalAskVolume := src.TotalAskVolume
	if math.IsNaN(totalAskVolume) {
		totalAskVolume = 0
	}

	totalBidVolume := src.TotalBidVolume
	if math.IsNaN(totalBidVolume) {
		totalBidVolume = 0
	}

	turnNum := src.TurnNum
	if math.IsNaN(turnNum) {
		turnNum = 0
	}

	dst := &Snapshot{
		RecordCircle:       src.RecordCircle,
		ExchangeTime:       src.ExchangeTime,
		TimestampS:         src.TimestampS,
		TimestampNS:        src.TimestampNS,
		InstrumentId:       ClearString(src.InstrumentId[:]),
		LastPrice:          src.LastPrice,
		TradingDay:         ClearString(src.TradingDay[:]),
		ExchangeId:         ClearString(src.ExchangeId[:]),
		OpenPrice:          src.OpenPrice,
		HighestPrice:       src.HighestPrice,
		LowestPrice:        src.LowestPrice,
		PreClosePrice:      src.PreClosePrice,
		Volume:             src.Volume,
		Amount:             src.Amount,
		ClosePrice:         src.ClosePrice,
		SettlementPrice:    src.SettlementPrice,
		PreSettlementPrice: src.PreSettlementPrice,
		OpenInterest:       src.OpenInterest,
		AveragePrice:       src.AveragePrice,
		AskPrice:           askPrice[:],
		AskVolume:          askVolume[:],
		BidPrice:           bidPrice[:],
		BidVolume:          bidVolume[:],
		AverageAskPrice:    averageAskPrice,
		AverageBidPrice:    averageBidPrice,
		TotalAskVolume:     totalAskVolume,
		TotalBidVolume:     totalBidVolume,
		UpdateTime:         ClearString(src.UpdateTime[:]),
		MilliSeconds:       src.MilliSeconds,
		UpperLimitPrice:    src.UpperLimitPrice,
		LowerLimitPrice:    src.LowerLimitPrice,
		Iopv:               src.Iopv,
		TurnNum:            turnNum,
	}

	return dst
}

// order 的自定义序列化

func CopyOrder(src *shm.Order) *Order {
	dst := &Order{
		Code:         src.Code,
		LocalTime:    src.LocalTime,
		ExchangeTime: src.ExchangeTime,
		Order:        src.Order,
		Price:        src.Price,
		Volume:       src.Volume,
		ChannelNo:    src.ChannelNo,
		OrderKind:    src.OrderKind,
		FunctionCode: src.FunctionCode,
		InstrumentId: ClearString(src.InstrumentId[:]),
		TradingDay:   ClearString(src.TradingDay[:]),
	}
	return dst
}

func CopyOrderExtra(src *shm.OrderExtra) *OrderExtra {
	dst := &OrderExtra{
		Order:    CopyOrder(&src.Order),
		OrderNo:  src.OrderNo,
		BizIndex: src.BizIndex,
	}
	return dst
}

// transaction 的自定义序列化

func CopyTransaction(src *shm.Transaction) *Transaction {
	dst := &Transaction{
		Code:         src.Code,
		LocalTime:    src.LocalTime,
		ExchangeTime: src.ExchangeTime,
		AskOrder:     src.AskOrder,
		BidOrder:     src.BidOrder,
		Price:        src.Price,
		Amount:       src.Amount,
		Volume:       src.Volume,
		ChannelNo:    src.ChannelNo,
		Index:        src.Index,
		BsFlag:       src.BsFlag,
		FunctionCode: src.FunctionCode,
		IsSzSe:       src.IsSzSe,
		InstrumentId: ClearString(src.InstrumentId[:]),
		TradingDay:   ClearString(src.TradingDay[:]),
		OrderKind:    src.OrderKind,
	}

	return dst
}

func CopyTransactionExtra(src *shm.TransactionExtra) *TransactionExtra {
	dst := &TransactionExtra{
		Transaction: CopyTransaction(&src.Transaction),
		BizIndex:    src.BizIndex,
	}
	return dst
}
