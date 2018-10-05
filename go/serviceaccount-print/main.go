package main

import (
	"io"
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		acc, err := appengine.ServiceAccount(ctx)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		io.WriteString(w, acc)
	})
}
