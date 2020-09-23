package options

import (
	"encoding/json"
	"bytes"
	"os"
	"testing"
	"io/ioutil"
)

func TestOptionChainSchema(t *testing.T) {
	//f, err := os.Open("testdata/spy_options_response.json")
	f, err := os.Open("testdata/nan_test.json")
	if err != nil {
		t.Error(err)
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	buf = bytes.Replace(buf, []byte("\"NaN\""), []byte("null"), -1)

	var chain OptionChain
	if err := json.NewDecoder(bytes.NewBuffer(buf)).Decode(&chain); err != nil {
		t.Error(err)
	}
}
