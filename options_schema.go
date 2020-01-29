package tdam

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type OptionChain struct {
	Symbol           string                       `json:"symbol"`
	Status           string                       `json:"status"`
	Underlying       *Underlying                  `json:"underlying"`
	Strategy         string                       `json:"strategy"`
	Interval         float64                      `json:"interval"`
	IsDelayed        bool                         `json:"isDelayed"`
	IsIndex          bool                         `json:"isIndex"`
	DaysToExpiration float64                      `json:"daysToExpiration"`
	InterestRate     float64                      `json:"interestRate"`
	UnderlyingPrice  float64                      `json:"underlyingPrice"`
	Volatility       float64                      `json:"volatility"`
	RawCalls         map[ExpirationDate]StrikeMap `json:"callExpDateMap"`
	RawPuts          map[ExpirationDate]StrikeMap `json:"putExpDateMap"`
}

func (c *OptionChain) ExpirationDates() []ExpirationDate {
	out := []ExpirationDate{}
	for exp := range c.RawCalls {
		out = append(out, exp)
	}
	sort.Sort(byDTE(out))
	return out
}

func (c *OptionChain) StrikeTable(exp ExpirationDate) StrikeTable {
	table := []Strike{}

	putMap := c.RawPuts[exp]
	callMap := c.RawCalls[exp]

	for price := range callMap {
		table = append(table, Strike{
			Price: price,
			Call:  callMap[price],
			Put:   putMap[price],
		})
	}
	sort.Sort(byPrice(table))
	return table
}

type Option struct {
	PutCall                string               `json:"putCall"`
	Symbol                 string               `json:"symbol"`
	Description            string               `json:"description"`
	ExchangeName           string               `json:"exchangeName"`
	BidPrice               float64              `json:"bid"`
	AskPrice               float64              `json:"ask"`
	LastPrice              float64              `json:"last"`
	MarkPrice              float64              `json:"mark"`
	BidSize                float64              `json:"bidSize"`
	AskSize                float64              `json:"askSize"`
	LastSize               float64              `json:"lastSize"`
	HighPrice              float64              `json:"highPrice"`
	LowPrice               float64              `json:"lowPrice"`
	OpenPrice              float64              `json:"openPrice"`
	ClosePrice             float64              `json:"closePrice"`
	TotalVolume            int64                `json:"totalVolume"`
	QuoteTimeInLong        float64              `json:"quoteTimeInLong"`
	TradeTimeInLong        float64              `json:"tradeTimeInLong"`
	NetChange              float64              `json:"netChange"`
	Volatility             float64              `json:"volatility"`
	Delta                  float64              `json:"delta"`
	Gamma                  float64              `json:"gamma"`
	Theta                  float64              `json:"theta"`
	Vega                   float64              `json:"vega"`
	Rho                    float64              `json:"rho"`
	TimeValue              float64              `json:"timeValue"`
	OpenInterest           float64              `json:"openInterest"`
	IsInTheMoney           bool                 `json:"isInTheMoney"`
	TheoreticalOptionValue float64              `json:"theoreticalOptionValue"`
	TheoreticalVolatility  float64              `json:"theoreticalVolatility"`
	IsMini                 bool                 `json:"isMini"`
	IsNonStandard          bool                 `json:"isNonStandard"`
	OptionDeliverablesList []OptionDeliverables `json:"optionDeliverablesList"`
	StrikePrice            float64              `json:"strikePrice"`
	ExpirationDate         string               `json:"expirationDate"`
	ExpirationType         string               `json:"expirationType"`
	Multiplier             float64              `json:"multiplier"`
	SettlementType         string               `json:"settlementType"`
	DeliverableNote        string               `json:"deliverableNote"`
	IsIndexOption          bool                 `json:"isIndexOption"`
	PercentChange          float64              `json:"percentChange"`
	MarkChange             float64              `json:"markChange"`
	MarkPercentChange      float64              `json:"markPercentChange"`
}

type StrikePrice float64

func (s StrikePrice) String() string {
	return fmt.Sprintf("%.1f", s)
}

func (v *StrikePrice) UnmarshalText(b []byte) error {
	f, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		return err
	}
	*v = StrikePrice(f)
	return nil
}

type Strike struct {
	Price StrikePrice
	Call  []Option // it's possible to get multiple exchanges' orderbook
	Put   []Option
}

type StrikeTable []Strike

func (s StrikeTable) NearestToDelta(d float64) (put, call Strike) {
	put = s[len(s)-1]
	call = s[0]

	for _, strike := range s {
		if strike.DistToDelta(d) < put.DistToDelta(d) {
			if math.Abs(strike.Put[0].Delta)-d < math.Abs(strike.Call[0].Delta)-d {
				put = strike
			}
		}
		if strike.DistToDelta(d) < call.DistToDelta(d) && strike.Call[0].Delta > 0 {
			if math.Abs(strike.Put[0].Delta)-d > math.Abs(strike.Call[0].Delta)-d {
				call = strike
			}
		}
	}

	return
}

func (s Strike) DistToDelta(d float64) float64 {
	c := s.Call[0]
	p := s.Put[0]
	return math.Abs(math.Min(math.Abs(c.Delta)-d, math.Abs(p.Delta)-d))
}

func (s Strike) DeltaAbove(d float64) bool {
	return math.Abs(s.Call[0].Delta) > d && math.Abs(s.Put[0].Delta) > d
}

func (s Strike) DeltaBelow(d float64) bool {
	return math.Abs(s.Call[0].Delta) < d || math.Abs(s.Put[0].Delta) < d
}

func (s Strike) DeltaBetween(a, b float64) bool {
	return s.DeltaAbove(a) && s.DeltaBelow(b)
}

type StrikeMap map[StrikePrice][]Option

type byPrice []Strike

func (s byPrice) Len() int {
	return len(s)
}

func (s byPrice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPrice) Less(i, j int) bool {
	return s[i].Price < s[j].Price
}

type byDTE []ExpirationDate

func (s byDTE) Len() int {
	return len(s)
}

func (s byDTE) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDTE) Less(i, j int) bool {
	return s[i].DTE() < s[j].DTE()
}

type ExpirationDate string

func (e ExpirationDate) String() string {
	fields := strings.Split(string(e), ":")
	return fmt.Sprintf("%s (%s)", fields[0], fields[1])
}

func (e ExpirationDate) Date() time.Time {
	fields := strings.Split(string(e), ":")
	date, err := time.Parse("2006-01-02", fields[0])
	if err != nil {
		fmt.Printf("invalid expiration date '%s': %v\n", e, err)
		date = time.Time{}
	}
	return date
}

func (e ExpirationDate) DTE() int {
	fields := strings.Split(string(e), ":")
	dte, err := strconv.ParseInt(fields[1], 10, 16)
	if err != nil {
		fmt.Printf("invalid expiration date '%s': %v\n", e, err)
		dte = -1
	}
	return int(dte)
}

func (ch *OptionChain) NearestDTE(target int) ExpirationDate {
	dates := ch.ExpirationDates()
	if len(dates) == 0 {
		return ExpirationDate("")
	}
	nearest := dates[0]
	for _, ed := range dates {
		if abs(target-ed.DTE()) < abs(target-nearest.DTE()) {
			nearest = ed
		}
	}
	return nearest
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Underlying struct {
	Ask               float64 `json:"ask"`
	AskSize           int     `json:"askSize"`
	Bid               float64 `json:"bid"`
	BidSize           float64 `json:"bidSize"`
	Change            float64 `json:"change"`
	Close             float64 `json:"close"`
	Delayed           bool    `json:"delayed"`
	Description       string  `json:"description"`
	ExchangeName      string  `json:"exchangeName"`
	FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
	FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
	HighPrice         float64 `json:"highPrice"`
	Last              float64 `json:"last"`
	LowPrice          float64 `json:"lowPrice"`
	Mark              float64 `json:"mark"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
	OpenPrice         float64 `json:"openPrice"`
	PercentChange     float64 `json:"percentChange"`
	QuoteTime         float64 `json:"quoteTime"`
	Symbol            string  `json:"symbol"`
	TotalVolume       int64   `json:"totalVolume"`
	TradeTime         int64   `json:"tradeTime"`
}

type OptionDeliverables struct {
	Symbol           string `json:"symbol"`
	AssetType        string `json:"assetType"`
	DeliverableUnits string `json:"deliverableUnits"`
	CurrencyType     string `json:"currencyType"`
}
