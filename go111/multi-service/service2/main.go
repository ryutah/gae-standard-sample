package main

import (
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/service2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello service2!!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
