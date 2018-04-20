package main

import (
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	r.Path("/projects/{project}/datasets/{dataset}/tables").Methods("GET").HandlerFunc(getDataset)

	http.Handle("/", r)
}

func getDataset(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	vars := mux.Vars(r)
	var (
		project = vars["project"]
		dataset = vars["dataset"]
	)

	cli, err := bigquery.NewClient(ctx, project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cli.Close()

	it := cli.Dataset(dataset).Tables(ctx)
	var tabls []string
	for {
		tbl, err := it.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tabls = append(tabls, tbl.TableID)
	}

	fmt.Fprintf(w, "%v", strings.Join(tabls, ", "))
}
