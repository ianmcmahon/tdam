package options

import (
	"encoding/json"
	"os"
	"testing"
)

func TestOptionChainSchema(t *testing.T) {
	//f, err := os.Open("testdata/spy_options_response.json")
	f, err := os.Open("testdata/880587956.json")
	if err != nil {
		t.Error(err)
	}

	var chain OptionChain
	if err := json.NewDecoder(f).Decode(&chain); err != nil {
		t.Error(err)
	}
}
