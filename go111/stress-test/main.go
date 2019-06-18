package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ret := fib(33)
		fmt.Fprintln(w, ret)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func fib(n uint64) int {
	switch n {
	case 0:
		return 0
	case 1, 2:
		return 1
	default:
		return fib(n-2) + fib(n-1)
	}
}
