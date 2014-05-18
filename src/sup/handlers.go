package sup

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func StartHTTPServer() {
	// Rest API
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/status/", statusHandler)

	// Resources
	http.HandleFunc("/", staticHandler) // html, js, css

	// Serve HTTP on localhost only. Let Nginx terminate HTTPS for us.
	address := fmt.Sprintf("127.0.0.1:%d", GetConfig().HTTPPort)
	log.Printf("Listening on http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, recoverAndLogHandler(http.DefaultServeMux)))
}

// Wraps an HTTP handler, adding error logging.
//
// If the inner function panics, the outer function recovers, logs, sends an
// HTTP 500 error response.
func recoverAndLogHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the ResponseWriter to remember the status
		rww := &responseWriterWrapper{0, w}
		begin := time.Now()

		defer func() {
			// Send a 500 error if a panic happens during a handler.
			// Without this, Chrome & Firefox were retrying aborted ajax requests,
			// at least to my localhost.
			if e := recover(); e != nil {
				rww.writeHeader(http.StatusInternalServerError)
                rww.Write([]byte(fmt.Sprintf("Error: %s", e)))
				log.Printf("%s: %s", e, debug.Stack())
			}

			// Finally, log.
			durationMS := time.Since(begin).Nanoseconds() / 1000000
			log.Printf("%s %s %v %v %s", r.RemoteAddr, r.Method, rww.Status, durationMS, r.URL)
		}()

		handler.ServeHTTP(rww, r)
	})
}

// Remember the status for logging
type responseWriterWrapper struct {
	Status int
	http.ResponseWriter
}

func (w *responseWriterWrapper) writeHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}
