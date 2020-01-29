package user

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/ianmcmahon/tdam"
)

type TokenTimestamp time.Time

var ttFormat string = "2006-01-02T15:04:05-0700"

func (t *TokenTimestamp) UnmarshalJSON(b []byte) error {
	tt, err := time.Parse(ttFormat, strings.Trim(string(b), `"`))
	if err != nil {
		return err
	}
	*t = TokenTimestamp(tt)
	return nil
}

type StreamerInfo struct {
	StreamerBinaryUrl string         `json:"streamerBinaryUrl"`
	StreamerSocketUrl string         `json:"streamerSocketUrl"`
	Token             string         `json:"token"`
	TokenTimestamp    TokenTimestamp `json:"tokenTimestamp"`
	UserGroup         string         `json:"userGroup"`
	AccessLevel       string         `json:"accessLevel"`
	Acl               string         `json:"acl"`
	AppId             string         `json:"appId"`
}

type QuotesDelayed struct {
	IsNyseDelayed   bool `json:"isNyseDelayed"`
	IsNasdaqDelayed bool `json:"isNasdaqDelayed"`
	IsOpraDelayed   bool `json:"isOpraDelayed"`
	IsAmexDelayed   bool `json:"isAmexDelayed"`
	IsCmeDelayed    bool `json:"isCmeDelayed"`
	IsIceDelayed    bool `json:"isIceDelayed"`
	IsForexDelayed  bool `json:"isForexDelayed"`
}

type SubscriptionKeys []string

func (s *SubscriptionKeys) UnmarshalJSON(b []byte) error {
	*s = []string{}
	var unrollMap map[string][]map[string]string
	if err := json.Unmarshal(b, &unrollMap); err != nil {
		return err
	}
	for _, inner := range unrollMap["keys"] {
		*s = append(*s, inner["key"])
	}
	return nil
}

type UPAccount struct {
	AccountId         string         `json:"accountId"`
	Description       string         `json:"description"`
	DisplayName       string         `json:"displayName"`
	AccountCdDomainId string         `json:"accountCdDomainId"`
	Company           string         `json:"company"`
	Segment           string         `json:"segment"`
	SurrogateIds      interface{}    `json:"surrogateIds"`
	Preferences       Preferences    `json:"preferences"`
	Acl               string         `json:"acl"`
	Authorizations    Authorizations `json:"authorizations"`
}

type Authorizations struct {
	Apex               bool   `json:"apex"`
	LevelTwoQuotes     bool   `json:"levelTwoQuotes"`
	StockTrading       bool   `json:"stockTrading"`
	MarginTrading      bool   `json:"marginTrading"`
	StreamingNews      bool   `json:"streamingNews"`
	OptionTradingLevel string `json:"optionTradingLevel"` //"'COVERED' or 'FULL' or 'LONG' or 'SPREAD' or 'NONE'",
	StreamerAccess     bool   `json:"streamerAccess"`
	AdvancedMargin     bool   `json:"advancedMargin"`
	ScottradeAccount   bool   `json:"scottradeAccount"`
}

type Preferences struct {
	ExpressTrading                   bool   `json:"expressTrading"`
	DirectOptionsRouting             bool   `json:"directOptionsRouting"`
	DirectEquityRouting              bool   `json:"directEquityRouting"`
	DefaultEquityOrderLegInstruction string `json:"defaultEquityOrderLegInstruction"` //"'BUY' or 'SELL' or 'BUYtOcOVER' or 'SELLsHORT' or 'NONE'",
	DefaultEquityOrderType           string `json:"defaultEquityOrderType"`           //"'MARKET' or 'LIMIT' or 'STOP' or 'STOPlIMIT' or 'TRAILINGsTOP' or 'MARKEToNcLOSE' or 'NONE'",
	DefaultEquityOrderPriceLinkType  string `json:"defaultEquityOrderPriceLinkType"`  //"'VALUE' or 'PERCENT' or 'NONE'",
	DefaultEquityOrderDuration       string `json:"defaultEquityOrderDuration"`       //"'DAY' or 'GOODtILLcANCEL' or 'NONE'",
	DefaultEquityOrderMarketSession  string `json:"defaultEquityOrderMarketSession"`  //"'AM' or 'PM' or 'NORMAL' or 'SEAMLESS' or 'NONE'",
	DefaultEquityQuantity            int    `json:"defaultEquityQuantity"`
	MutualFundTaxLotMethod           string `json:"mutualFundTaxLotMethod"`    //"'FIFO' or 'LIFO' or 'HIGHcOST' or 'LOWcOST' or 'MINIMUMtAX' or 'AVERAGEcOST' or 'NONE'",
	OptionTaxLotMethod               string `json:"optionTaxLotMethod"`        //"'FIFO' or 'LIFO' or 'HIGHcOST' or 'LOWcOST' or 'MINIMUMtAX' or 'AVERAGEcOST' or 'NONE'",
	EquityTaxLotMethod               string `json:"equityTaxLotMethod"`        //"'FIFO' or 'LIFO' or 'HIGHcOST' or 'LOWcOST' or 'MINIMUMtAX' or 'AVERAGEcOST' or 'NONE'",
	DefaultAdvancedToolLaunch        string `json:"defaultAdvancedToolLaunch"` //"'TA' or 'N' or 'Y' or 'TOS' or 'NONE' or 'CC2'",
	AuthTokenTimeout                 string `json:"authTokenTimeout"`          //"'FIFTYfIVEmINUTES' or 'TWO_HOURS' or 'FOUR_HOURS' or 'EIGHT_HOURS'"
}

type UserPrincipal struct {
	AuthToken                string            `json:"authToken"`
	UserId                   string            `json:"userId"`
	UserCdDomainId           string            `json:"userCdDomainId"`
	PrimaryAccountId         string            `json:"primaryAccountId"`
	LastLoginTime            string            `json:"lastLoginTime"`
	TokenExpirationTime      string            `json:"tokenExpirationTime"`
	LoginTime                string            `json:"loginTime"`
	AccessLevel              string            `json:"accessLevel"`
	StalePassword            bool              `json:"stalePassword"`
	StreamerInfo             StreamerInfo      `json:"streamerInfo"`
	ProfessionalStatus       string            `json:"professionalStatus"` // "'PROFESSIONAL' or 'NON_PROFESSIONAL' or 'UNKNOWN_STATUS'",
	Quotes                   QuotesDelayed     `json:"quotes"`
	StreamerSubscriptionKeys *SubscriptionKeys `json:"streamerSubscriptionKeys"`
	Accounts                 []UPAccount       `json:"accounts"`
}

func GetUserPrincipals(c *tdam.Client) (*UserPrincipal, error) {
	token, err := c.TDAMToken()
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: &tls.Config{}}
	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("GET", "https://api.tdameritrade.com/v1/userprincipals", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	options := req.URL.Query()
	options.Set("fields", "streamerSubscriptionKeys,streamerConnectionInfo")
	//options.Add("fields", "streamerConnectionInfo")
	req.URL.RawQuery = options.Encode()

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

	var u UserPrincipal
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}
