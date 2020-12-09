package tdam

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type TDTime time.Time

func (e TDTime) UnmarshalJSON(b []byte) (err error) {
	s := string(b)

	// Get rid of the quotes "" around the value.
	// A second option would be to include them
	// in the date format string instead, like so below:
	//   time.Parse(`"`+time.RFC3339Nano+`"`, s)
	s = s[1 : len(s)-1]

	t, err := time.Parse("2006-01-02T15:04:05-0700", s)
	if err == nil {
		e = TDTime(t)
	}
	return
}

type Expiration TDTime

func (e Expiration) UnmarshalJSON(b []byte) (err error) {
	t := TDTime(e)
	err = t.UnmarshalJSON(b)
	if err == nil {
		e = Expiration(t)
	}
	return
}

func (e Expiration) String() string {
	return time.Time(e).Format("02 Jan 2006")
}

type Transaction struct {
	Type                          string           `json:"type"`
	ClearingReferenceNumber       string           `json:"clearingReferenceNumber"`
	SubAccount                    string           `json:"subAccount"`
	SettlementDate                string           `json:"settlementDate"`
	OrderID                       string           `json:"orderId"`
	SMA                           float64          `json:"sma"`
	RequirementReallocationAmount float64          `json:"requirementReallocationAmount"`
	DayTradeBPE                   float64          `json:"dayTradeBuyingPowerEffect"`
	NetAmount                     float64          `json:"netAmount"`
	TransactionDate               TDTime           `json:"transactionDate"`
	OrderDate                     TDTime           `json:"orderDate"`
	TransactionSubType            string           `json:"transactionSubType"`
	TransactionID                 int              `json:"transactionId"`
	CashBalanceEffectFlag         bool             `json:"cashBalanceEffectFlag"`
	Description                   string           `json:"description"`
	ACHStatus                     string           `json:"achStatus"`
	AccruedInterest               float64          `json:"accruedInterest"`
	Fees                          interface{}      `json:"fees"`
	TransactionItem               *TransactionItem `json:"transactionItem"`
}

type TransactionItem struct {
	AccountID            int         `json:"accountId"`
	Amount               float64     `json:"amount"`
	Price                float64     `json:"price"`
	Cost                 float64     `json:"cost"`
	ParentOrderKey       int         `json:"parentOrderKey"`
	ParentChildIndicator string      `json:"parentChildIndicator"`
	Instruction          string      `json:"instruction"`
	PositionEffect       string      `json:"positionEffect"`
	Instrument           *Instrument `json:"instrument"`
}

type Instrument struct {
	AssetType            InstrumentType `json:"assetType"`
	Symbol               Symbol         `json:"symbol"`
	UnderlyingSymbol     Symbol         `json:"underlyingSymbol"`
	OptionExpirationDate Expiration     `json:"optionExpirationDate"`
	OptionStrikePrice    float64        `json:"optionStrikePrice"`
	PutCall              string         `json:"putCall"`
	CUSIP                string         `json:"cusip"`
	Description          string         `json:"description"`
	BondMaturityDate     string         `json:"bondMaturityDate"`
	BondInterestRate     string         `json:"bondInterestRate"`
}

func (t *Transaction) String() string {
	i := t.TransactionItem
	return fmt.Sprintf("%s to %s %.0f @ %.2f: %#v", i.Instruction, i.PositionEffect, i.Amount, i.Price, i.Instrument)
}

func (i *Instrument) String() string {
	switch i.AssetType {
	case OPTION:
		return fmt.Sprintf("%s %s %g %s", i.UnderlyingSymbol,
			i.OptionExpirationDate, i.OptionStrikePrice, i.PutCall)
	default:
		return fmt.Sprintf("Unknown asset type: %s", i.AssetType)
	}
}

func (i *Instrument) populateFromSymbol() error {
	if i.AssetType != OPTION {
		return nil
	}

	re := regexp.MustCompile(`([^_]+)_(\d{6})(P|C)(\d+(.\d+)?)`)
	matches := re.FindStringSubmatch(string(i.Symbol))

	if len(matches) < 6 {
		return fmt.Errorf("couldn't parse option symbol")
	}

	if v, err := strconv.ParseFloat(matches[4], 32); err != nil {
		return err
	} else {
		i.OptionStrikePrice = v
	}

	if t, err := time.Parse("010206", matches[2]); err != nil {
		return err
	} else {
		i.OptionExpirationDate = Expiration(t)
	}

	return nil
}
