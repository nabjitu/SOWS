package main

import (
	//"encoding/json"
	"fmt"
	//"net/http"
	//"strconv"

	//"google.golang.org/appengine"
	//"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/urlfetch"
	"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/log"
	"time"
	//"cloud.google.com/go/storage"
	//"golang.org/x/net/context"
	//"cloud.google.com/go/datastore"
	"net/http"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type Post struct {
	Title       string
	Body        string `datastore:",noindex"`
	PublishedAt time.Time
}

type Entity struct {
	Value string
}

func handle(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	k := datastore.NewKey(ctx, "Entity", "stringID", 0, nil)
	e := new(Entity)
	if err := datastore.Get(ctx, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	old := e.Value
	e.Value = r.URL.Path

	if _, err := datastore.Put(ctx, k, e); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "old=%q\nnew=%q\n", old, e.Value)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// first create a new context
	c := appengine.NewContext(r)
	// and use that context to create a new http client
	client := urlfetch.Client(c)

	// now we can use that http client as before
	res, err := client.Get("http://google.com")
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hallo Got Google with status %s\n", res.Status)
	fmt.Fprintf(w, "\n Det virker her er koden %s\n", res.Status)

	fmt.Fprintf(w, "\n Hey HEY.. Cool %s\n", res.Status)
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/handle", handle)
}


func main() {

	// ScopeDatastore grants permissions to view and/or manage datastore entities
	//const ScopeDatastore = "https://www.googleapis.com/auth/datastore"
	//Vores
	//const ScopeDatastore = "http://storage.googleapis.com/[BUCKET_NAME]/[OBJECT_NAME]"

	url := "http://storage.googleapis.com/"
	BUCKET_NAME := ""
	OBJECT_NAME := ""
	ScopeDatastore := url + BUCKET_NAME + OBJECT_NAME
	fmt.Print(ScopeDatastore)

	http.HandleFunc("/", handler)
	http.HandleFunc("/handle", handle)
}

