package loghelper

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/IAmRDhar/scaling-web-app/backend/entity"
)

var logserviceURL = flag.String("logservice", "http://172.18.0.14:6000", "Address of logging service")

var tr = http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

var client = &http.Client{Transport: &tr}

func WriteEntry(entry *entity.LogEntry) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	enc.Encode(entry)
	req, _ := http.NewRequest(http.MethodPost, *logserviceURL, &buf)
	fmt.Printf("Sending log request (%v) to %s\n", &buf, *logserviceURL)
	client.Do(req)
}
