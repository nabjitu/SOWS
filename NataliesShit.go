package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
)

//Product bla bla
type Product struct {
	Name        string `json:"name"`
	Supermarket string `json:"supermarket"`
	Price       int64  `json:"price"`
}

func printProduct(w http.ResponseWriter, r *http.Request, p Product) {
	fmt.Fprintf(w, "%v , ", p.Name)
	fmt.Fprintf(w, "%v , ", p.Supermarket)
	fmt.Fprintf(w, "%v", p.Price)
}

//Add a new item to the viewAll(w, r)
//Virker!
func addHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Old viewAll(w, r) contained: %v \n", viewAll(w, r))
	firstValue := r.FormValue("Name")
	secondValue := r.FormValue("Supermarket")
	thirdValue, err := strconv.ParseInt(r.FormValue("Price")[0:], 10, 64)
	if err != nil {
		fmt.Fprint(w, "Price has to be integer \n")
	}
	//-------Datastore-------
	// create a new App Engine context from the HTTP request.
	ctx := appengine.NewContext(r)

	p := &Product{Name: firstValue, Supermarket: secondValue, Price: thirdValue}

	// create a new complete key of kind Person and value gopher.
	key := datastore.NewKey(ctx, "Product", firstValue, 0, nil)
	// put p in the datastore.
	key, err = datastore.Put(ctx, key, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Fprintf(w, "%v was NOT stored!", key)
		return
	}
	fmt.Fprintf(w, "%v was stored i datastore!", key)
}

//Virker
func viewAll(w http.ResponseWriter, r *http.Request) (p []Product) {
	ctx := appengine.NewContext(r)
	q := datastore.NewQuery("Product")
	q = q.Order("Name")

	// and finally execute the query retrieving all values into p.
	_, err := q.GetAll(ctx, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return p
}

func getAll(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var p []Product
	q := datastore.NewQuery("Product")
	q = q.Order("Name")
	_, err := q.GetAll(ctx, &p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, p)
}

func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var p Product

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Product is %v exist in %v for %v kr.", p.Name, p.Supermarket, p.Price)
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/jsons", decodeHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// first create a new context
	c := appengine.NewContext(r)
	// and use that context to create a new http client
	client := urlfetch.Client(c)

	//Lav url
	url := "http://storage.googleapis.com/"
	bucketName := "gcp-public-data-sentinel-2"
	objectName := "/tiles/01/C/CV/S2A_MSIL1C_20151221T205519_N0201_R028_T01CCV_20160329T181515.SAFE/GRANULE/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_N02.01/IMG_DATA/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_B02.jp2"
	ScopeDatastore := url + bucketName + objectName

	// now we can use that http client as before
	res, err := client.Get(ScopeDatastore)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Hallo Got Google with status %s\n", res.Status)
}

/*func mux() {
	mux := http.NewServeMux()
	mux.Handler("/", handler)
	http.ListenAndServe(":9000", nil)
}
*/
func main() {
	fmt.Print("woop woopp")
	mux := http.NewServeMux()
	//files := http.FileServer(http.Dir(config.Static)) mux.Handle("/static/", http.StripPrefix("/static/", files))
	//mux.HandleFunc("/", index)

	mux.HandleFunc("/handler", handler)
	//http.HandleFunc("/add", addHandler)
	//http.HandleFunc("/jsons", decodeHandler)
}