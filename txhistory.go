package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (a *Account) TradeHistory(symbol Symbol) ([]Transaction, error) {
	token, err := a.TDAMToken()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	endpoint := fmt.Sprintf("https://api.tdameritrade.com/v1/accounts/%s/transactions", a.AccountId)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	query := req.URL.Query()
	query.Add("type", "TRADE")
	query.Add("symbol", string(symbol))
	query.Add("startDate", "2019-10-01")
	query.Add("endDate", time.Now().Format("2006-01-02"))
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

	/*
		dump, err := httputil.DumpResponse(resp, true)
		fmt.Printf("%s\n", dump)
	*/

	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
