package streamer

import (
	"strconv"
	"strings"
	"time"
)

type requestWrapper map[string][]request

type request struct {
	Service    string            `json:"service"`
	Command    string            `json:"command"`
	RequestID  int               `json:"requestid"`
	Account    string            `json:"account"`
	Source     string            `json:"source"`
	Parameters map[string]string `json:"parameters"`
}

func (r request) isLoginCommand() bool {
	return r.Service == "ADMIN" && r.Command == "LOGIN"
}

type responseWrapper struct {
	Response []response `json:"response"`
	Notify   []notify   `json:"notify"`
	Data     []Data     `json:"data"`
	//Snapshot []snapshot `json:"snapshot"`
}

type notify struct {
	Heartbeat wsTimestamp `json:"heartbeat"`
}

type response struct {
	Service   string                 `json:"service"`
	Command   string                 `json:"command"`
	Timestamp wsTimestamp            `json:"timestamp"`
	RequestID int                    `json:"requestid,string"`
	Content   map[string]interface{} `json:"content"`
}

type Data struct {
	Service   string                   `json:"service"`
	Command   string                   `json:"command"`
	Timestamp wsTimestamp              `json:"timestamp"`
	Content   []map[string]interface{} `json:"content"`
}

type wsTimestamp time.Time

func (t *wsTimestamp) UnmarshalJSON(b []byte) error {
	millis, err := strconv.ParseInt(strings.Trim(string(b), `"`), 10, 64)
	if err != nil {
		return err
	}
	secs := millis / 1000
	millis = millis - (secs * 1000)
	*t = wsTimestamp(time.Unix(secs, millis*1000))

	return nil
}
