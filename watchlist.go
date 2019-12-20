package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Watchlist struct {
	Name        string          `json:"name"`
	WatchlistId string          `json:"watchlistId"`
	AccountId   string          `json:"accountId"`
	Status      string          `json:"status"`
	Items       []WatchlistItem `json:"watchlistItems"`
}

type WatchlistItem struct {
	SequenceId   int                 `json:"sequenceId"`
	Quantity     float64             `json:"quantity"`
	AveragePrice float64             `json:"averagePrice"`
	Commission   float64             `json:"commission"`
	PurchaseDate string              `json:"purchasedDate"`
	Instrument   WatchlistInstrument `json:"instrument"`
	Status       string              `json:"status"`
}

type WatchlistInstrument struct {
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	AssetType   string `json:"assetType"`
}

func GetWatchlists() ([]Watchlist, error) {
	token, err := TDAMToken()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/accounts/watchlists", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

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

	var watchlists []Watchlist
	if err := json.NewDecoder(resp.Body).Decode(&watchlists); err != nil {
		return nil, err
	}

	return watchlists, nil
}
