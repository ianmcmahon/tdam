package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Scanner struct {
	Authenticated bool
}

func DTE(min, max time.Duration) (from, to string) {
	now := time.Now()
	near := now.Add(min * time.Hour * 24)
	far := now.Add(max * time.Hour * 24)

	return near.Format("2006-01-02"), far.Format("2006-01-02")
}

func (ch *OptionChain) divideTable(table StrikeTable) (puts, calls StrikeTable) {
	for i, strike := range table {
		if float64(strike.Price) > ch.UnderlyingPrice {
			puts, calls = table[:i], table[i:]
			//sort.Reverse(byPrice(puts))
			return
		}
	}
	return
}

func (s *Scanner) Coinflip(symbol string, strikeWidth int) (*OptionChain, error) {
	// filter by RSI
	// choose directional bias based on crossingbelow30 bearish or crossingabove70 for bullish
	// for now we'll choose bullish
	// if we're bullish, we sell an ATM put spread, collecting a large premium as a bet that
	// the stock will go up
	//direction := PUT

	fromDate, toDate := DTE(20, 45)

	options := url.Values{
		//"range":    []string{"NTM"}, // Near The Money
		"strategy":    []string{"SINGLE"},
		"strikeCount": []string{"100"},
		"fromDate":    []string{fromDate},
		"toDate":      []string{toDate},
	}

	chain, err := s.GetChain(symbol, options)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	exp := chain.NearestDTE(45)
	//fmt.Printf("targeting expiration %s\n", exp)
	puts, _ := chain.divideTable(chain.StrikeTable(exp))
	//fmt.Printf("underlying price: %.2f\n", chain.UnderlyingPrice)

	if len(puts) > 1 {
		shortStrike := puts[len(puts)-1]
		longStrike := puts[len(puts)-1-strikeWidth]
		spreadWidth := float64(shortStrike.Price - longStrike.Price)
		credit := shortStrike.Put[0].BidPrice - longStrike.Put[0].BidPrice
		if credit/spreadWidth > 0.5 {
			fmt.Printf("%s %s PUT %.1f | %.2fΔ cr %.2f/%.1f (%.1f)%%\n", symbol, exp,
				shortStrike.Price, shortStrike.Put[0].Delta, credit, spreadWidth,
				credit/spreadWidth*100.0)
		}
	}

	return chain, nil
}

func (s *Scanner) GetChain(symbol string, options url.Values) (*OptionChain, error) {
	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/marketdata/chains", nil)
	if err != nil {
		return nil, err
	}
	if options == nil {
		options = req.URL.Query()
	}
	if s.Authenticated {
		token, err := TDAMToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		options.Set("apikey", ConsumerKey)
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
