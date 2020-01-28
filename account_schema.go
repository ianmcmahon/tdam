package tdam

import "fmt"

type AccountType string
type InstrumentType string
type Symbol string

const (
	CASH   AccountType = "CASH"
	MARGIN             = "MARGIN"

	OPTION          InstrumentType = "OPTION"
	MUTUAL_FUND                    = "MUTUAL_FUND"
	CASH_EQUIVALENT                = "CASH_EQUIVALENT"
	EQUITY                         = "EQUITY"
	INDEX                          = "INDEX"
	FIXED_INCOME                   = "FIXED_INCOME"
	CURRENCY                       = "CURRENCY"
)

type Account struct {
	*Client
	SecuritiesAccount `json:"securitiesAccount"`
}

type SecuritiesAccount struct {
	Type                    AccountType `json:"type"`
	AccountId               string      `json:"accountId"`
	RoundTrips              int         `json:"roundTrips"`
	IsDayTrader             bool        `json:"isDayTrader"`
	IsClosingOnlyRestricted bool        `json:"isClosingOnlyRestricted"`
	RawPositions            []*Position `json:"positions"`
	InitialBalances         Balances    `json:"initialBalances"`
	CurrentBalances         Balances    `json:"currentBalances"`
	ProjectedBalances       Balances    `json:"projectedBalances"`
	//OrderStrategies         []OrderStrategy `json:"orderStrategies"`
}

func (a SecuritiesAccount) Positions() map[Symbol][]*Position {
	out := make(map[Symbol][]*Position)
	for _, p := range a.RawPositions {
		sym := p.Instrument.UnderlyingSymbol
		if _, ok := out[sym]; !ok {
			out[sym] = make([]*Position, 0)
		}
		out[sym] = append(out[sym], p)
	}
	return out
}

type Position struct {
	ShortQuantity                  float64     `json:"shortQuantity"`
	AveragePrice                   float64     `json:"averagePrice"`
	CurrentDayProfitLoss           float64     `json:"currentDayProfitLoss"`
	CurrentDayProfitLossPercentage float64     `json:"currentDayProfitLossPercentage"`
	LongQuantity                   float64     `json:"longQuantity"`
	SettledLongQuantity            float64     `json:"settledLongQuantity"`
	SettledShortQuantity           float64     `json:"settledShortQuantity"`
	AgedQuantity                   float64     `json:"agedQuantity"`
	Instrument                     *Instrument `json:"instrument"`
	MarketValue                    float64     `json:"marketValue"`
}

func (p Position) String() string {
	return fmt.Sprintf("%g %s %.2f", p.LongQuantity-p.ShortQuantity, p.Instrument, p.MarketValue)
}

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
