package datatype

const (
	OrderFuncBuy  = "B"
	OrderFuncSell = "S"
	OrderFunKnown = "C"
)

const (
	OrderKindMkt = "1" // 市价
	OrderKindFix = "2" // 限价

	OrderKindUsf = "U" // 本方最优
	OrderKindUcf = "Y" // 对手方最优
	//OrderKindUtp = "2" // 即时成交
)

const (
	TransactionFuncCancel = "C" // 撤单
	TransactionFuncTrans  = "F" // 成交
)

const (
	TransactionBSFlagBuy     = "B"
	TransactionBSFlagSell    = "S"
	TransactionBSFlagUnknown = "N"
)

const (
	TypeSnapshot              = 0
	TypeOrder                 = 1
	TypeTransaction           = 2
	TypeOrderTransaction      = 10
	TypeSnapshotExtra         = 11
	TypeOrderExtra            = 12
	TypeTransactionExtra      = 13
	TypeOrderTransactionExtra = 14
)

// 正常类型

/*
snapshot: {"RecordCircle":188064729037483404,"ExchangeTime":91927000,"TimestampS":1690161567,"TimestampNS":131970697,"InstrumentId":"000101.SH","LastPrice":229.4614,"TradingDay":"20230724","ExchangeId":"","OpenPrice":229.4614,"HighestPrice":229.4614,"LowestPrice":229.4614,"PreClosePrice":229.4001,"Volume":0,"Amount":0,"ClosePrice":0,"SettlementPrice":0,"PreSettlementPrice":0,"OpenInterest":0,"AveragePrice":0,"AskPrice":[0,0,0,0,0,0,0,0,0,0],"AskVolume":[0,0,0,0,0,0,0,0,0,0],"BidPrice":[0,0,0,0,0,0,0,0,0,0],"BidVolume":[0,0,0,0,0,0,0,0,0,0],"AverageAskPrice":0,"AverageBidPrice":0,"TotalAskVolume":0,"TotalBidVolume":0,"UpdateTime":"09:19:27","MilliSeconds":0,"UpperLimitPrice":0,"LowerLimitPrice":0,"Iopv":0,"TurnNum":0}
*/
type Snapshot struct {
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

/*
{"Code":2317,"LocalTime":1690161885658230008,"ExchangeTime":92445590,"Order":358031,"Price":17,"Volume":700,"ChannelNo":2013,"OrderKind":"2","FunctionCode":"S","InstrumentId":"002317.SZ","TradingDay":"20230724","OrderNo":0,"BizIndex":0}
*/
type Order struct {
	Code         uint64
	LocalTime    uint64 // 1690161885658230008
	ExchangeTime uint64 // 92445590
	Order        int64
	Price        float64
	Volume       uint32

	ChannelNo    int32
	OrderKind    string // 1=市价单   2=限价单   U=best_passive  3=IOC; K=ETF; V=best_5_levels; W=FOK; X=best_passive; Y=best_aggressive
	FunctionCode string // buy='B'   sell='S'  no 'C'

	InstrumentId string
	TradingDay   string
}

/*
transactionExtra: {"Code":2400,"LocalTime":1690162201247357587,"ExchangeTime":93000090,"AskOrder":468957,"BidOrder":439218,"Price":6.05,"Amount":1210,"Volume":200,"ChannelNo":2013,"Index":468958,"BsFlag":"S","FunctionCode":"F","IsSzSe":true,"InstrumentId":"002400.SZ","TradingDay":"20230724","OrderKind":48,"BizIndex":0}
*/
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
	BsFlag       string // buy='B'   sell='S'    unknown='N'
	FunctionCode string // fill='F'  cancel='C'
	IsSzSe       bool
	InstrumentId string
	TradingDay   string
	OrderKind    byte // 0
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
