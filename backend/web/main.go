package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IAmRDhar/scaling-web-app/backend/util"
	"github.com/IAmRDhar/scaling-web-app/backend/web/controller"
)

var loadbalancerURL = flag.String("loadbalancer", "https://172.18.0.12:2001", "Address of the load balancer")

func main() {
	flag.Parse()

	templateCache, _ := buildTemplateCache()
	controller.Setup(templateCache)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	go http.ListenAndServeTLS(":3000", "/cert.pem", "/key.pem", new(util.GzipHandler))

	go func() {
		for range time.Tick(300 * time.Millisecond) {
			tc, isUpdated := buildTemplateCache()
			if isUpdated {
				controller.SetTemplateCache(tc)
			}
		}
	}()

	// Could have done using channels
	// @todo TODO see waitgroups
	time.Sleep(2 * time.Second)

	// Making a service discovery request
	// as soon as the app is up
	http.Get(*loadbalancerURL + "/register?port=3000")

	log.Println("Server started, press <ENTER> to exit")
	fmt.Scanln()

	// Inform Load balancer that app is down
	http.Get(*loadbalancerURL + "/unregister?port=3000")

}

var lastModTime time.Time = time.Unix(0, 0)

func buildTemplateCache() (*template.Template, bool) {
	needUpdate := false

	f, err := os.Open("/web/templates")
	if err != nil {
		//	panic(err)
	}
	fileInfos, _ := f.Readdir(-1)
	fileNames := make([]string, len(fileInfos))
	for idx, fi := range fileInfos {
		if fi.ModTime().After(lastModTime) {
			lastModTime = fi.ModTime()
			needUpdate = true
		}
		fileNames[idx] = "/web/templates/" + fi.Name()
	}

	var tc *template.Template
	if needUpdate {
		log.Print("Template change detected, updating...")
		tc = template.Must(template.New("").ParseFiles(fileNames...))
		log.Println("template update complete")
	}
	return tc, needUpdate
}
