package data

import (
	"crypto/tls"
	"flag"
	"net/http"
)

var dataServiceUrl = flag.String("dataservice", "https://localhost:4000", "Address of the data service provider")

func init() {
	tr := http.Transport{
		// New transport config
		TLSClientConfig: &tls.Config{
			// In development we have and unsiged
			// certificate by default and therefore Go is gonna
			// have problems.
			// skipping verification of insecure HTTP ceritficates
			// Not be done in production
			InsecureSkipVerify: true,
		},
	}

	http.DefaultClient = &http.Client{Transport: &tr}
}
