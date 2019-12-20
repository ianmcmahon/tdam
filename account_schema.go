package tdam

type AccountType string
type InstrumentType string

const (
	CASH   AccountType = "CASH"
	MARGIN             = "MARGIN"

	OPTION          InstrumentType = "Option"
	MUTUAL_FUND                    = "MutualFund"
	CASH_EQUIVALENT                = "CashEquivalent"
	EQUITY                         = "Equity"
	FIXED_INCOME                   = "FixedIncome"
)

type Account struct {
	SecuritiesAccount `json:"securitiesAccount"`
}

type SecuritiesAccount struct {
	Type                    AccountType `json:"type"`
	AccountId               string      `json:"accountId"`
	RoundTrips              int         `json:"roundTrips"`
	IsDayTrader             bool        `json:"isDayTrader"`
	IsClosingOnlyRestricted bool        `json:"isClosingOnlyRestricted"`
	Positions               []Position  `json:"positions"`
	InitialBalances         Balances    `json:"initialBalances"`
	CurrentBalances         Balances    `json:"currentBalances"`
	ProjectedBalances       Balances    `json:"projectedBalances"`
	//OrderStrategies         []OrderStrategy `json:"orderStrategies"`
}

type Position struct {
	ShortQuantity                  float64    `json:"shortQuantity"`
	AveragePrice                   float64    `json:"averagePrice"`
	CurrentDayProfitLoss           float64    `json:"currentDayProfitLoss"`
	CurrentDayProfitLossPercentage float64    `json:"currentDayProfitLossPercentage"`
	LongQuantity                   float64    `json:"longQuantity"`
	SettledLongQuantity            float64    `json:"settledLongQuantity"`
	SettledShortQuantity           float64    `json:"settledShortQuantity"`
	AgedQuantity                   float64    `json:"agedQuantity"`
	Instrument                     Instrument `json:"instrument"`
	MarketValue                    float64    `json:"marketValue"`
}

type Instrument interface{}

type Balances struct {
	AccruedInterest              float64 `json:"accruedInterest"`
	CashBalance                  float64 `json:"cashBalance"`
	CashReceipts                 float64 `json:"cashReceipts"`
	LongOptionMarketValue        float64 `json:"longOptionMarketValue"`
	LiquidationValue             float64 `json:"liquidationValue"`
	LongMarketValue              float64 `json:"longMarketValue"`
	MoneyMarketFund              float64 `json:"moneyMarketFund"`
	Savings                      float64 `json:"savings"`
	ShortMarketValue             float64 `json:"shortMarketValue"`
	PendingDeposits              float64 `json:"pendingDeposits"`
	CashAvailableForTrading      float64 `json:"cashAvailableForTrading"`
	CashAvailableForWithdrawal   float64 `json:"cashAvailableForWithdrawal"`
	CashCall                     float64 `json:"cashCall"`
	LongNonMarginableMarketValue float64 `json:"longNNMarginableMarketValue"`
	TotalCash                    float64 `json:"totalCash"`
	ShortOptionMarketValue       float64 `json:"shortOptionMarketValue"`
	MutualFundValue              float64 `json:"mutualFundValue"`
	BondValue                    float64 `json:"bondValue"`
	CashDebitCallValue           float64 `json:"cashDebitCallValue"`
	UnsettledCash                float64 `json:"unsettledCash"`
}
