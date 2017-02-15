package eintel

import (
	"bytes"
	"encoding/json"
	"errors"
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

func JumpCount(from, to string) int {
	url := fmt.Sprintf("http://api.eve-central.com/api/route/from/%s/to/%s", from, to)
	route := make([]interface{}, 0)

	if response, err := getCached(url); err == nil {
		json.NewDecoder(response).Decode(&route)
		return len(route)
	}

	return 0

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
	return nil, errors.New("No Response")
}
