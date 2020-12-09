package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (a *Account) TradeHistoryCallback(symbol Symbol, cb func(symbol Symbol, data []byte)) error {
	token, err := a.TDAMToken()
	if err != nil {
		return err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	endpoint := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/transactions", a.AccountId)

	current := time.Now()
	for {
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		query := req.URL.Query()
		query.Add("type", "TRADE")
		query.Add("symbol", string(symbol))
		query.Add("startDate", current.AddDate(0, -1, 0).Format("2006-01-02"))
		query.Add("endDate", current.Format("2006-01-02"))
		req.URL.RawQuery = query.Encode()
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("status %d: %s\n", resp.StatusCode, resp.Status)
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		val := []interface{}{}
		if err := json.Unmarshal(body, &val); err != nil {
			return err
		}
		if len(val) == 0 {
			return nil
		}

		cb(symbol, body)

		current = current.AddDate(0, -1, 0)
	}

	return nil
}

func (a *Account) TradeHistory(symbol Symbol) (out []Transaction, err error) {
	out = []Transaction{}
	a.TradeHistoryCallback(symbol, func(symbol Symbol, data []byte) {
		var transactions []Transaction
		//fmt.Printf("%s: %s\n", symbol, data)
		if e := json.Unmarshal(data, &transactions); e != nil {
			err = e
		} else {
			out = append(out, transactions...)
		}
	})

	return
}
