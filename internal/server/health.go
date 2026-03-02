package server

import (
	"net/http"
)

func StartHealthServer() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	go http.ListenAndServe(":8080", nil)
}