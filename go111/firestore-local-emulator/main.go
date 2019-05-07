// Package main provides ...
package main

import (
	"net/http"

	"fmt"
	"math/rand"

	"os"

	"encoding/json"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	port      = os.Getenv("PORT")
)

// Foo is foo
type Foo struct {
	Name string
}

func main() {
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Close()

		newFoo := Foo{
			Name: fmt.Sprintf("Foo%v", rand.Int()),
		}
		ref, _, err := client.Collection("foo").Add(ctx, newFoo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":"%v"}`, ref.ID)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := firestore.NewClient(ctx, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Close()

		docs := client.Collection("foo").Documents(ctx)
		foos := make([]Foo, 0)
		for {
			ref, err := docs.Next()
			if err == iterator.Done {
				break
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var foo Foo
			if err := ref.DataTo(&foo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			foos = append(foos, foo)
		}
		byt, err := json.MarshalIndent(foos, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(byt)
	})

	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
