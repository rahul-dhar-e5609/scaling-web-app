package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type cacheEntry struct {
	data       []byte
	expiration time.Time
}

var (
	cache = make(map[string]*cacheEntry)
	mutex = sync.RWMutex{}
)

var maxAgeRegexp = regexp.MustCompile(`maxage=(\d+)`)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getFromCache(w, r)
		} else if r.Method == http.MethodPost {
			saveToCache(w, r)
		}
	})

	http.HandleFunc("/invalidate", invalidateEntry)

	go http.ListenAndServeTLS(":5000", "/cert.pem", "/key.pem", nil)

	log.Println("Caching service started, press <ENTER> to exit")

	fmt.Scanln()
}

func invalidateEntry(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()
	key := r.URL.Query().Get("key")
	fmt.Printf("Purging entry with key %s\n", key)
	delete(cache, key)
}

func getFromCache(w http.ResponseWriter, r *http.Request) {
	mutex.RLock()
	defer mutex.RUnlock()

	key := r.URL.Query().Get("key")
	fmt.Printf("Searching cahce for %s...\n", key)
	if entry, ok := cache[key]; ok {
		fmt.Println("Found")
		w.Write(entry.data)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Println("Not found")
}

func saveToCache(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	key := r.URL.Query().Get("key")
	cacheHeader := r.Header.Get("cache-control")

	fmt.Printf("Saving cache entry with key %s for %s seconds\n", key, cacheHeader)

	matches := maxAgeRegexp.FindStringSubmatch(cacheHeader)
	if len(matches) == 2 {
		dur, _ := strconv.Atoi(matches[1])
		data, _ := ioutil.ReadAll(r.Body)
		cache[key] = &cacheEntry{data: data, expiration: time.Now().Add(time.Duration(dur) * time.Second)}
	}
}
