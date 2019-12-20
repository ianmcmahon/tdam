package tdam

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

func (c *Client) GetAccounts() ([]Account, error) {
	token, err := c.TDAMToken()
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
