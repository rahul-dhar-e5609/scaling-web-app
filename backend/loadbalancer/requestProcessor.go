package main

import (
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/IAmRDhar/scaling-web-app/backend/entity"
	"github.com/IAmRDhar/scaling-web-app/backend/logservice/loghelper"
)

var (
	// URL of all the application servers
	appservers = []string{}
	// Which application to call / is called
	currentIndex = 0
	// to forward the request to those app servers
	client = http.Client{Transport: &transport}
)

func processRequests() {
	for {
		select {
		case request := <-requestCh:
			println("request")
			if len(appservers) == 0 {
				request.w.WriteHeader(http.StatusInternalServerError)
				request.w.Write([]byte("No app srevers found"))
				request.doneCh <- struct{}{}
				continue
			}
			currentIndex++
			if currentIndex == len(appservers) {
				currentIndex = 0
			}
			host := appservers[currentIndex]
			go processRequest(host, request)
		case host := <-registerCh:
			println("register: " + host)
			go loghelper.WriteEntry(&entity.LogEntry{
				Level:     entity.LogLevelInfo,
				Timestamp: time.Now(),
				Source:    "load balancer",
				Message:   "Registering application server with address: " + host,
			})
			isFound := false
			for _, h := range appservers {
				if host == h {
					isFound = true
					break
				}
			}
			if !isFound {
				appservers = append(appservers, host)
			}
		case host := <-unregisterCh:
			println("unregister: " + host)
			go loghelper.WriteEntry(&entity.LogEntry{
				Level:     entity.LogLevelInfo,
				Timestamp: time.Now(),
				Source:    "load balancer",
				Message:   "Unregistering application server with address: " + host,
			})
			for i := len(appservers) - 1; i >= 0; i-- {
				if appservers[i] == host {
					appservers = append(appservers[:i], appservers[i+1:]...)
				}
			}
		case <-heartbeatCh:
			println("heartbeat")
			// Copy the appservers slice so that it isnt
			// affected while the appserver thinning is done
			server := appservers[:]

			// Every request to check if the appservers
			// are up or not is done in this go routine,
			// mutually exclusive to other channels, no
			// memory / appservers slice is being shared
			go func(servers []string) {
				for _, h := range servers {
					resp, err := http.Get("https://" + h + "/ping")
					if err != nil || resp.StatusCode != 200 {
						// unregistering
						unregisterCh <- h
					}
				}
			}(server)
		}
	}
}

func processRequest(host string, request *webRequest) {
	hostURL, _ := url.Parse(request.r.URL.String())
	hostURL.Scheme = "https"
	hostURL.Host = host
	println(host)
	println(hostURL.String())
	req, _ := http.NewRequest(request.r.Method, hostURL.String(), request.r.Body)
	for k, v := range request.r.Header {
		//	values := ""
		for _, headerValue := range v {
			// println("[Request to ", host, "] Header: ", k, " Value: ", headerValue)
			//		values += headerValue + " "
			req.Header.Add(k, headerValue)
		}
		// req.Header.Add(k, values)
	}

	resp, err := client.Do(req)

	if err != nil {
		request.w.WriteHeader(http.StatusInternalServerError)
		request.doneCh <- struct{}{}
		return
	}

	for k, v := range resp.Header {
		//values := ""
		for _, headerValue := range v {
			// println("[Response from ", host, "] Header: ", k, " Value: ", headerValue)
			//	values += headerValue + " "
			request.w.Header().Add(k, headerValue)
		}
		// request.w.Header().Add(k, values)
	}
	request.w.Header().Add("Content-Type", "text/html")
	io.Copy(request.w, resp.Body)

	request.doneCh <- struct{}{}
}
