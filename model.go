package md

const (
	TypeMarketData            = 0
	TypeOrder                 = 1
	TypeTransaction           = 2
	TypeOrderTransaction      = 10
	TypeMarketDataExtra       = 11
	TypeOrderExtra            = 12
	TypeTransactionExtra      = 13
	TypeOrderTransactionExtra = 14
)

// 正常类型

type MarketData struct {
	RecordCircle       uint64
	ExchangeTime       int64
	TimestampS         int64
	TimestampNS        int64
	InstrumentId       string
	LastPrice          float64
	TradingDay         string
	ExchangeId         string
	OpenPrice          float64
	HighestPrice       float64
	LowestPrice        float64
	PreClosePrice      float64
	Volume             float64
	Amount             float64
	ClosePrice         float64
	SettlementPrice    float64
	PreSettlementPrice float64
	OpenInterest       float64
	AveragePrice       float64
	AskPrice           []float64
	AskVolume          []float64
	BidPrice           []float64
	BidVolume          []float64
	AverageAskPrice    float64
	AverageBidPrice    float64
	TotalAskVolume     float64
	TotalBidVolume     float64
	UpdateTime         string
	MilliSeconds       int32
	UpperLimitPrice    float64
	LowerLimitPrice    float64
	Iopv               float64
	TurnNum            float64
}

type Order struct {
	Code         uint64
	LocalTime    uint64
	ExchangeTime uint64
	Order        int64
	Price        float64
	Volume       uint32

	ChannelNo    int32
	OrderKind    byte
	FunctionCode byte

	InstrumentId string
	TradingDay   string
}

type Transaction struct {
	Code         uint64
	LocalTime    uint64
	ExchangeTime uint64
	AskOrder     int64
	BidOrder     int64
	Price        float64
	Amount       uint64
	Volume       uint32

	ChannelNo    int32
	Index        int32
	BsFlag       byte
	FunctionCode byte
	IsSzSe       bool
	InstrumentId string
	TradingDay   string
	OrderKind    byte
}

// extra 类型

type OrderExtra struct {
	*Order

	OrderNo  int64
	BizIndex int64
}

type TransactionExtra struct {
	*Transaction
	BizIndex int64
}
