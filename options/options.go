package options

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"io/ioutil"
  "bytes"

	"github.com/ianmcmahon/tdam"
)

type Client struct {
	*tdam.Client
}

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

func (c *Client) GetChain(symbol string, options url.Values) (*OptionChain, error) {
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
		token, err := c.TDAMToken()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	} else {
		options.Set("apikey", c.ConsumerKey)
	}
	options.Set("symbol", symbol)
	req.URL.RawQuery = options.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Printf("status %d: %s\n", resp.StatusCode, resp.Status)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var chain OptionChain
	if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(&chain); err != nil {
		tmpdir, _ := ioutil.TempDir("/tmp", "tdam-debug-*")
		tmpfile, _ := ioutil.TempFile(tmpdir, "*.json")
		n, errr := tmpfile.Write(body)
		if errr != nil {
			fmt.Printf("error writing dumpfile: %v\n", errr)
		}
		fmt.Printf("logged (%d bytes) bad json to %s\n", n, tmpfile.Name())
		return nil, err
	}

	return &chain, nil
}
