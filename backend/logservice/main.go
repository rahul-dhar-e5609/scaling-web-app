package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/IAmRDhar/scaling-web-app/backend/entity"
)

var mutex sync.Mutex
var entries logEntries

const logPath = "/log/log.txt"

var tickCh = time.Tick(5 * time.Second)
var writeDelay = 5 * time.Second

func main() {
	flag.Parse()

	http.HandleFunc("/", storeEntry)

	f, _ := os.Create(logPath)
	f.Close()

	go http.ListenAndServeTLS(":6000", "./cert.pem", "./key.pem", nil)
	go writeLog()

	log.Println("Log service started, press <ENTER> to exit")
	fmt.Scanln()
}

func storeEntry(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var entry entity.LogEntry
	err := dec.Decode(&entry)

	fmt.Println(entry, err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	mutex.Lock()
	entries = append(entries, entry)
	mutex.Unlock()
}

func writeLog() {
	for range tickCh {
		mutex.Lock()

		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0664)
		if err != nil {
			fmt.Println(err)
			continue
		}

		targetTime := time.Now().Add(-writeDelay)
		sort.Sort(entries)
		for i, entry := range entries {
			if entry.Timestamp.Before(targetTime) {
				_, err := logFile.WriteString(writeEntry(entry))
				if err != nil {
					fmt.Println(err)
				}
				if i == len(entries)-1 {
					entries = logEntries{}
				}
			} else {
				entries = entries[i:]
				break
			}
		}

		logFile.Close()

		mutex.Unlock()
	}
}

func writeEntry(entry entity.LogEntry) string {
	return fmt.Sprintf("%v;%v;%v;%v\n",
		entry.Timestamp.Format("2006-01-02 15:04:05"),
		entry.Level, entry.Source, entry.Message)
}
