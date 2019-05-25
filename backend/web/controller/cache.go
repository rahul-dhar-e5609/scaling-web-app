package controller

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var cacheServiceURL = flag.String("cachingservice", "https://172.18.0.13:5000", "Address of the caching service")

func getFromCache(key string) (io.ReadCloser, bool) {
	resp, err := http.Get(*cacheServiceURL + "/?key=" + key)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("Get failed with response code %v | ERROR: %v\n", resp.StatusCode, err)
		return nil, false
	}
	return resp.Body, true
}

func savingToCache(key string, duration int64, data []byte) {
	req, _ := http.NewRequest(http.MethodPost, *cacheServiceURL+"/?key="+key, bytes.NewBuffer(data))
	req.Header.Add("cache-control", "maxage="+strconv.FormatInt(duration, 10))
	http.DefaultClient.Do(req)
}

func invalideCacheEntry(key string) {
	http.Get(*cacheServiceURL + "/invalidate?key=" + key)
}
