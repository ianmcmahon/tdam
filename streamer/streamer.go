package streamer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ianmcmahon/tdam"
	"github.com/ianmcmahon/tdam/user"
)

const (
	QOSExpress  QoSLevel = "0" // 500 ms
	QOSRealTime QoSLevel = "1" // 750 ms
	QOSFast     QoSLevel = "2" // 1000ms
	QOSModerate QoSLevel = "3" // 1500ms
	QOSSlow     QoSLevel = "4" // 3000 ms
	QOSDelayed  QoSLevel = "5" // 3000 ms

)

type QoSLevel string
type responseCallback func(resp response)
type DataCallback func(symbol string, resp Data)

type Streamer struct {
	tdamClient   *tdam.Client
	principal    *user.UserPrincipal
	conn         *websocket.Conn
	done         chan bool
	requestCount int
	wg           *sync.WaitGroup

	responseCallbacks map[int]responseCallback // maps by requestID

	// maps by-->    service     symbol    subscriber
	dataCallbacks map[string]map[string]map[string]DataCallback
	cbMutex       *sync.RWMutex

	// maps by-->  service     symbol  subscriber
	subscribers map[string]map[string][]string
}

func (s *Streamer) nextRequest() int {
	cnt := s.requestCount
	s.requestCount++
	return cnt
}

func New(client *tdam.Client) (*Streamer, error) {
	up, err := user.GetUserPrincipals(client)
	if err != nil {
		return nil, err
	}
	s := &Streamer{
		tdamClient:        client,
		principal:         up,
		done:              make(chan bool),
		requestCount:      0,
		wg:                &sync.WaitGroup{},
		cbMutex:           &sync.RWMutex{},
		responseCallbacks: make(map[int]responseCallback),
		dataCallbacks: map[string]map[string]map[string]DataCallback{
			"QUOTE":                    make(map[string]map[string]DataCallback),
			"OPTION":                   make(map[string]map[string]DataCallback),
			"LEVELONE_FUTURES":         make(map[string]map[string]DataCallback),
			"LEVELONE_FUTURES_OPTIONS": make(map[string]map[string]DataCallback),
		},
		subscribers: map[string]map[string][]string{
			"QUOTE":                    make(map[string][]string),
			"OPTION":                   make(map[string][]string),
			"LEVELONE_FUTURES":         make(map[string][]string),
			"LEVELONE_FUTURES_OPTIONS": make(map[string][]string),
		},
	}

	// keep wg held whenever we aren't logged in
	s.wg.Add(1)

	return s, nil
}

func (s *Streamer) Run() error {
	var err error

	if s == nil {
		return fmt.Errorf("streamer is nil!?")
	}

	if s.principal == nil {
		return fmt.Errorf("No user principal!")
	}

	u := url.URL{Scheme: "ws", Host: s.principal.StreamerInfo.StreamerSocketUrl, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	s.conn, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	go s.handleIncoming()

	if err := s.sendRequest(s.loginRequest(), func(resp response) {
		code := resp.Content["code"]

		// successful login releases the waitlock
		if v, ok := code.(float64); ok && v == 0.0 {
			//fmt.Printf("Login: code: %T %v\n", code, code)
			s.wg.Done()
		}
	}); err != nil {
		log.Printf("error sending login request: %v\n", err)
		return err
	}

	return nil
}

func (s *Streamer) QoS(level QoSLevel) error {

	if err := s.sendRequest(s.qosRequest(level), func(resp response) {
	}); err != nil {
		log.Printf("error sending qos request: %v\n", err)
		return err
	}
	return nil
}

func (s *Streamer) qosRequest(level QoSLevel) request {

	loginReq := request{
		Service:   "ADMIN",
		Command:   "QOS",
		RequestID: s.nextRequest(),
		Account:   s.principal.Accounts[0].AccountId,
		Source:    s.principal.StreamerInfo.AppId,
		Parameters: map[string]string{
			"qoslevel": string(level),
		},
	}

	return loginReq
}

func (s *Streamer) Stop() error {

	if err := s.sendRequest(s.logoutRequest(), func(resp response) {
		code := resp.Content["code"]

		if v, ok := code.(float64); ok && v == 0.0 {
			//fmt.Printf("Logout: code: %T %v\n", code, code)
			s.done <- true
		}
	}); err != nil {
		log.Printf("error sending logout request: %v\n", err)
		return err
	}

	return nil
}

func (s *Streamer) logoutRequest() request {

	logoutReq := request{
		Service:    "ADMIN",
		Command:    "LOGOUT",
		RequestID:  s.nextRequest(),
		Account:    s.principal.Accounts[0].AccountId,
		Source:     s.principal.StreamerInfo.AppId,
		Parameters: map[string]string{},
	}

	return logoutReq
}

func (s *Streamer) handleIncoming() {
	//dump, err := os.Create("streamdump.log")
	//if err != nil {
	//		log.Fatal(err)
	//}
	//defer dump.Close()
	for {
		_, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			s.done <- true
			return
		}

		var resp responseWrapper
		if err := json.Unmarshal(message, &resp); err != nil {
			log.Printf("received: %s", message)
			log.Printf("error unmarshaling response: %v", err)
			continue
		}

		// there's three (four really) possible responses in here
		// is it possible that we get different types of response in one packet?
		// maybe!  let's just spin through all of them and route the messages
		//fmt.Printf("recv: %#v\n", resp)
		for _, resp := range resp.Response {
			// requests are submitted with a requestID, and the caller
			// can attach a callback.  When we receive the response for that request,
			// route it back to the caller
			if cb, ok := s.responseCallbacks[resp.RequestID]; ok {
				cb(resp)
			}
		}
		for _, data := range resp.Data {
			// data messages come from subscriptions.
			// need to find the list of subscribers for this service/symbol and notify all
			for _, packet := range data.Content {
				symbol, ok := packet["key"].(string)
				if !ok {
					continue
				}
				//fmt.Printf("dataCallbacks[%s][%s: %v\n", data.Service, symbol, s.dataCallbacks[data.Service][symbol])
				s.cbMutex.RLock()
				callbacks := s.dataCallbacks[data.Service][symbol]
				s.cbMutex.RUnlock()
				for _, cb := range callbacks {
					cb(symbol, Data{data.Service, data.Command, data.Timestamp, []map[string]interface{}{packet}})
				}
			}
		}
	}
}

func (s *Streamer) sendRequest(req request, cf responseCallback) error {
	if !req.isLoginCommand() {
		s.wg.Wait()
	}

	// todo: this is hacky and could be better.  Assumes only one request in the wrapper.
	requestId := req.RequestID

	wrappedReq := requestWrapper{"requests": []request{req}}

	s.responseCallbacks[requestId] = cf
	data, err := json.Marshal(wrappedReq)
	if err != nil {
		return err
	}

	//log.Printf("sending req: %s\n", data)
	if err := s.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("write:", err)
		return err
	}

	return nil
}

func (s *Streamer) loginRequest() request {
	timestamp := time.Time(s.principal.StreamerInfo.TokenTimestamp).Unix() * 1000

	credentials := url.Values{
		"userid":      []string{s.principal.Accounts[0].AccountId},
		"token":       []string{s.principal.StreamerInfo.Token},
		"company":     []string{s.principal.Accounts[0].Company},
		"segment":     []string{s.principal.Accounts[0].Segment},
		"cddomain":    []string{s.principal.Accounts[0].AccountCdDomainId},
		"usergroup":   []string{s.principal.StreamerInfo.UserGroup},
		"accesslevel": []string{s.principal.StreamerInfo.AccessLevel},
		"authorized":  []string{"Y"},
		"timestamp":   []string{fmt.Sprintf("%d", timestamp)},
		"appid":       []string{s.principal.StreamerInfo.AppId},
		"acl":         []string{s.principal.StreamerInfo.Acl},
	}

	loginReq := request{
		Service:   "ADMIN",
		Command:   "LOGIN",
		RequestID: s.nextRequest(),
		Account:   s.principal.Accounts[0].AccountId,
		Source:    s.principal.StreamerInfo.AppId,
		Parameters: map[string]string{
			"credential": url.QueryEscape(credentials.Encode()),
			"token":      s.principal.StreamerInfo.Token,
			"version":    "1.0",
		},
	}

	return loginReq
}
