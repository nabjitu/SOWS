package main

import (
	//n"encoding/json"
	"fmt"
	"net/http"
	//n"strconv"

	//n"google.golang.org/appengine"
	//n"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/urlfetch"
	"io/ioutil"

	// Imports the Google Cloud BigQuery client package.
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"reflect"
)

//Add a new item to the viewAll(w, r)
//Virker!
/*func addHandler(w http.ResponseWriter, r *http.Request) {

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
}*/

//Virker
/*func viewAll(w http.ResponseWriter, r *http.Request) (p []Item) {
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
}*/


/*func decodeHandler(w http.ResponseWriter, r *http.Request) {
	var p Product

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Product is %v exist in %v for %v kr.", p.Name, p.Supermarket, p.Price)
}*/

func main() {
	//mux := http.NewServeMux()
	//files := http.FileServer(http.Dir(config.Static)) mux.Handle("/static/", http.StripPrefix("/static/", files))
	//mux.HandleFunc("/", index)

	http.HandleFunc("/", handler)
	//http.HandleFunc("/jsons", decodeHandler)
	http.ListenAndServe(":9000", nil)
	http.HandleFunc("/a", askBigQuery)
	http.HandleFunc("/mgrs", makeMGRS)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// first create a new context
	//c := appengine.NewContext(r)
	// and use that context to create a new http client
	client := http.DefaultClient

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
	fmt.Println(ioutil.ReadAll(res.Body))
}

func askBigQuery(w http.ResponseWriter, r *http.Request) {
	firstValue := r.FormValue("northLatitude")
	//secondValue := r.FormValue("southLatitute")
	//thirdValue:= r.FormValue("westLongditude")
	//fourthValue := r.FormValue("eastLongditude")
	// and use that context to create a new http client
	//client := http.DefaultClient

	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, "bigquery-public-data:cloud_storage_geo_index.sentinel_2_index")
	if err != nil {
		fmt.Fprintf(w," no go")
	}

	query := client.Query(
		`SELECT *
			FROM [bigquery-public-data:cloud_storage_geo_index.sentinel_2_index]
			WHERE north_lat =` + firstValue + `LIMIT 1000`)


/*
AND south_lat = ` + secondValue +`
			AND west_lon = ` + thirdValue +`
			AND east_lon = ` + fourthValue +`
*/
	// Execute the query.
	it, err := query.Read(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	// Iterate through the results.
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Println(values)
	}
	/*//Lav url
	url := "http://storage.googleapis.com/"
	bucketName := "gcp-public-data-sentinel-2"
	objectName := "/tiles/01/C/CV/S2A_MSIL1C_20151221T205519_N0201_R028_T01CCV_20160329T181515.SAFE/GRANULE/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_N02.01/IMG_DATA/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_B02.jp2"
	ScopeDatastore := url + bucketName + objectName

	//lav query
	//q :=

	// now we can use that http client as before
	res, err := client.Get(ScopeDatastore)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(ioutil.ReadAll(res.Body))*/
}

func useMgrsApi(w http.ResponseWriter, r *http.Request) (api string){
	firstValue := r.FormValue("northLatitude")
	//secondValue := r.FormValue("southLatitute")
	thirdValue:= r.FormValue("westLongditude")
	//fourthValue := r.FormValue("eastLongditude")

	//http://legallandconverter.com/cgi-bin/android5c.cgi?username=DEVELOPX&password=TEST1234&latitude=48.00820&longitude=-112.61440&cmd=mgrsrev1
	api = "http://legallandconverter.com/cgi-bin/android5c.cgi?username=DEVELOPX&password=TEST1234&latitude=" + firstValue + "&longitude=" + thirdValue + "&cmd=mgrsrev1"
	return
}

func makeMGRS(w http.ResponseWriter, r *http.Request) /*(result string)*/ {
	client := http.DefaultClient
	res, err := client.Get(useMgrsApi(w, r))
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println(reflect.TypeOf(res))
}
