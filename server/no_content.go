package server

import "net/http"

func noContentStatus(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}
