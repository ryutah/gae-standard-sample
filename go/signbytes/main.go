package main

import (
	"fmt"
	"io"
	"net/http"

	"google.golang.org/appengine"
)

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		_, sig, err := appengine.SignBytes(ctx, []byte("Hello world"))
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		io.WriteString(w, fmt.Sprintf("%s", sig))
	})
}
