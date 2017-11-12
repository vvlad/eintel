package eintel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GitbookIO/diskache"
	"io/ioutil"
	"net/http"
	// "os"
)

var (
	opts = &diskache.Opts{
		Directory: "/tmp/eintel/cache",
	}
	cache, cacheError = diskache.New(opts)
)

type KaelRoutes struct {
	Count int `json:"count"`
}

const (
	UnknownNumberOfJumps   = -1982
	IrelevantNumberOfJumps = -999
)

func JumpCount(from, to string) int {
	url := fmt.Sprintf("http://everest.kaelspencer.com/route/%s/%s/", from, to)
	route := KaelRoutes{}

	if response, err := getCached(url); err == nil {
		json.NewDecoder(response).Decode(&route)
		return route.Count
	}

	return UnknownNumberOfJumps
}

func getCached(uri string) (*bytes.Reader, error) {
	key := fmt.Sprintf("http:%s", uri)
	if data, inCache := cache.Get(key); inCache {
		return bytes.NewReader(data), nil
	} else {
		resp, err := http.Get(uri)
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		cache.Set(key, data)
		return bytes.NewReader(data), nil
	}
}
