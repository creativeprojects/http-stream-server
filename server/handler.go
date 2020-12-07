package server

import (
	"io"
	"net/http"
	"strings"

	"github.com/creativeprojects/clog"
)

func handler(w http.ResponseWriter, r *http.Request) {
	path := cleanupPath(r.URL.Path)
	clog.Debugf("%s %s", r.Method, path)

	switch r.Method {
	case "GET":
		sr := getTransfer(path)
		if !sr.lockReceiver() {
			// cannot get the same path twice
			w.WriteHeader(http.StatusConflict)
			break
		}
		sr.receiverChan <- w
		// wait ...
		select {
		// ... until the finish message
		case <-sr.senderFinishedChan:
			clog.Trace("receiver got sender finished message")
		// ... or a request cancellation
		case <-r.Context().Done():
			clog.Trace("receiver cancelled context")
			close(sr.receiverFinishedChan)
		}
	case "PUT":
		sr := getTransfer(path)
		select {
		case receiver := <-sr.receiverChan:
			receiver.Header().Add("Content-Type", "application/octet-stream")
			io.Copy(receiver, r.Body)
			close(sr.senderFinishedChan)
		case <-sr.receiverFinishedChan:
			clog.Trace("sender got receiver finished message")
		case <-r.Context().Done():
			clog.Trace("sender cancelled context")
		}
		deleteTransfer(path)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	clog.Debugf("%s %s finished", r.Method, path)
}

func cleanupPath(path string) string {
	return strings.TrimSuffix(strings.TrimSpace(path), "/")
}
