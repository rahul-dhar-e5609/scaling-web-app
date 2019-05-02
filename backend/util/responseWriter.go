package util

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// CloseableResponseWriter implements GZIP
// gZIP is a stream bsed protocol and
// we need to know when the stream ends
// so we can flush its buffer
//
// Interface because we need multiple
// responsewriters, eg. if the browser
// does not support gzip then we cant really use
// Gzip encoding with that request
type CloseableResponseWriter interface {
	http.ResponseWriter
	// Used for flushing the buffer
	Close()
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*gzip.Writer
}

func (w gzipResponseWriter) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (w gzipResponseWriter) Close() {
	// Closes the zgip writer
	w.Writer.Close()
}

func (w gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

type closeableResponseWriter struct {
	http.ResponseWriter
}

func (w closeableResponseWriter) Close() {

}

func GetResponseWriter(w http.ResponseWriter, req *http.Request) CloseableResponseWriter {
	// Check if browser can use the gzip encryption
	if strings.Contains(req.Header.Get("Account-Encoding"), "gzip") {
		w.Header().Set("Content-Encoding", "gzip")
		gRn := gzipResponseWriter{
			ResponseWriter: w,
			Writer:         gzip.NewWriter(w),
		}
		return gRn
	} else {
		return closeableResponseWriter{ResponseWriter: w}
	}
}

type GzipHandler struct{}

func (h *GzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseWriter := GetResponseWriter(w, r)
	defer responseWriter.Close()

	http.DefaultServeMux.ServeHTTP(responseWriter, r)
}
