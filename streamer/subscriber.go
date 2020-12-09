package streamer

import (
	"fmt"
	"log"
	"strings"
)

func (s *Streamer) Subscribe(service string, subscriber string, symbols []string, cb DataCallback) error {
	// filter out any symbols already subscribed to
	filteredSymbols := symbols[:0]
	for _, symbol := range symbols {
		// add our callback to it regardless
		s.cbMutex.Lock()
		if _, ok := s.dataCallbacks[service][symbol]; !ok {
			s.dataCallbacks[service][symbol] = make(map[string]DataCallback)
		}
		if _, ok := s.dataCallbacks[service][symbol][subscriber]; ok {
			return fmt.Errorf("'%s' already subscribed to %s/%s.  Use a unique subscriber name!", subscriber, service, symbol)
		}
		s.dataCallbacks[service][symbol][subscriber] = cb
		s.cbMutex.Unlock()

		if !s.isSubscribed(service, symbol, subscriber) {
			filteredSymbols = append(filteredSymbols, symbol)
			s.subscribers[service][symbol] = append(s.subscribers[service][symbol], subscriber)
		}
	}

	if err := s.sendRequest(s.subRequest(service, filteredSymbols), func(resp response) {
		//fmt.Printf("sub registration callback: %v\n", resp)
	}); err != nil {
		log.Printf("error sending subscribe request: %v\n", err)
		return err
	}

	return nil
}

func (s *Streamer) isSubscribed(service, symbol, subscriber string) bool {
	if subs, ok := s.subscribers[service][symbol]; ok {
		for _, sub := range subs {
			if sub == subscriber {
				return true
			}
		}
	}
	return false
}

func (s *Streamer) subRequest(service string, symbols []string) request {
	req := request{
		Service:   service,
		Command:   "SUBS",
		RequestID: s.nextRequest(),
		Account:   s.principal.Accounts[0].AccountId,
		Source:    s.principal.StreamerInfo.AppId,
		Parameters: map[string]string{
			"keys":   strings.Join(symbols, ","),
			"fields": "0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,49,51,52",
		},
	}

	return req
}
