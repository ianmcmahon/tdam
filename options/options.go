package options

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ianmcmahon/tdam"
)

func dte(min, max int) (from, to string) {
	now := time.Now()
	near := now.Add(time.Duration(min) * time.Hour * 24)
	far := now.Add(time.Duration(max) * time.Hour * 24)

	return near.Format("2006-01-02"), far.Format("2006-01-02")
}

func Single(strikes, fromDTE, toDTE int) url.Values {
	fromDate, toDate := dte(fromDTE, toDTE)

	options := url.Values{
		"strategy":    []string{"SINGLE"},
		"strikeCount": []string{fmt.Sprintf("%d", strikes)},
		"fromDate":    []string{fromDate},
		"toDate":      []string{toDate},
	}

	return options
}

func GetChain(symbol string, options url.Values) (*OptionChain, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/marketdata/chains", nil)
	if err != nil {
		return nil, err
	}
	if options == nil {
		options = req.URL.Query()
	}
	if false { //s.Authenticated {
		token, err := tdam.TDAMToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		options.Set("apikey", tdam.ConsumerKey)
	}
	options.Set("symbol", symbol)
	req.URL.RawQuery = options.Encode()

	//dump, _ := httputil.DumpRequest(req, false)
	//fmt.Println(string(dump))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("status %d: %s\n", resp.StatusCode, resp.Status)
	}

	var chain OptionChain
	if err := json.NewDecoder(resp.Body).Decode(&chain); err != nil {
		return nil, err
	}

	return &chain, nil
}
