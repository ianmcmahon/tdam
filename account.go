package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

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

func GetAccounts() ([]Account, error) {
	token, err := TDAMToken()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/accounts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	query := req.URL.Query()
	//query.Add("fields", "positions")
	query.Add("fields", "orders")
	req.URL.RawQuery = query.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("status %d: %s\n", resp.StatusCode, resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	fmt.Printf("%s\n", dump)

	var accounts []Account
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, err
	}

	return accounts, nil
}
