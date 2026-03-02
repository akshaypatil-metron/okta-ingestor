package server

import (
	"log"
	"net/http"
)

func StartHealthCheck() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	log.Println("Starting health check server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Health check server failed: %v", err)
	}
}
