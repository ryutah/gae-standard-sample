// Package main provides ...
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Response From k8s Service!!")
	})

	http.ListenAndServe(":8080", nil)
}
