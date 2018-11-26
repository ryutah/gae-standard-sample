package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/golang/glog"
)

func main() {
	flag.Set("stderrthreshold", "INFO")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		glog.Infof("hello!!")
		glog.Infof("context: %#v", r.Context())
		for k, v := range r.Header {
			glog.Infof("%s: %v", k, v)
		}
		for _, v := range os.Environ() {
			glog.Infof("%s", v)
		}
		fmt.Fprintln(w, "Hello World!!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
