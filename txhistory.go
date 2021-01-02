package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (a *Account) TradeHistoryCallback(symbol Symbol, start *time.Time, cb func(symbol Symbol, data []byte)) error {
	token, err := a.TDAMToken()
	if err != nil {
		return err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	endpoint := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/transactions", a.AccountId)

	current := time.Now()
	if current.Sub(*start) < 1*time.Hour {
		fmt.Printf("%s less than an hour old, skipping\n", symbol)
		return nil
	}
	for {
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

		prevMonth := current.AddDate(0, -6, 0)
		startTime := start
		if startTime == nil || startTime.Before(prevMonth) {
			startTime = &prevMonth
		}
		//fmt.Printf("%s: %s - %s\n", symbol, startTime, current)

		query := req.URL.Query()
		query.Add("type", "TRADE")
		query.Add("symbol", string(symbol))
		query.Add("startDate", startTime.Format("2006-01-02"))
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
		if len(val) == 0 && start == nil {
			return nil
		}

		cb(symbol, body)

		current = current.AddDate(0, -6, 0)
		if start != nil && current.Before(*start) {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (a *Account) TradeHistory(symbol Symbol, start *time.Time) (out []Transaction, err error) {
	out = []Transaction{}
	a.TradeHistoryCallback(symbol, start, func(symbol Symbol, data []byte) {
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
