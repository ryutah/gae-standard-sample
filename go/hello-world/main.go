package main

import (
	"net/http"
	"os"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/", helloworld)
	http.HandleFunc("/env", printEnv)
	http.HandleFunc("/appid", printAppID)
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!!"))
}

func printEnv(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "%v", os.Environ())
	w.Write([]byte("Print environ, see Stackdriver logging"))
}

func printAppID(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "%q", appengine.AppID(ctx))
	w.Write([]byte("Print AppID, see Stackdriver logging"))
}
