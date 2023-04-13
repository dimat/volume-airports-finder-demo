package service

import (
	"encoding/json"
	"net/http"
)

// Endpoint is a generic interface for a service endpoint.
type Endpoint[Req any, Resp any] interface {
	Call(Req) (Resp, error)
}

// HandlerFunc converts a generic endpoint to a http.HandlerFunc that can be used to serve
func HandlerFunc[Req any, Resp any](endpoint Endpoint[Req, Resp]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var args Req
		err := json.NewDecoder(r.Body).Decode(&args)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		reply, err := endpoint.Call(args)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(reply)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}
