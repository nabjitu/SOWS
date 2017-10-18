package main

import (
	//"encoding/json"
	"fmt"
	"net/http"
	//"strconv"
	//"os"

	//"google.golang.org/appengine"
	//"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/urlfetch"
	"io/ioutil"

	// Imports the Google Cloud BigQuery client package.
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	//"golang.org/x/net/context"
	//"golang.org/x/oauth2/google"
	//"google.golang.org/api/compute/v1"
	"github.com/im7mortal/UTM"
	
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
	http.HandleFunc("/a", askBigQuery)
	//http.HandleFunc("/jsons", decodeHandler)
	http.ListenAndServe(":9000", nil)
	
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

func makeMGRS(lat float, long float) {
	
	Mgrs_latbands := `CDEFGHIJKLMNPQRSTUVWXX`

	Mgrs_e100k := [`ABCDEFGH`, `JKLMNPQR`, `STUVWXYZ`]

	Mgrs_n100k := [`ABCDEFGHJKLMNPQRSTUV`, `FGHJKLMNPQRSTUVABCDE`]

	latLon := UTM.LatLon{lat, long}
	
		result, err := latLon.FromLatLon()
		if err != nil {
			panic(err.Error())
		}
		
			/*fmt.Printf(
				"Easting: %f; Northing: %f; ZoneNumber: %d; ZoneLetter: %s;",
				result.Easting,
				result.Northing,
				result.ZoneNumber,
				result.ZoneLetter,
			)*/
		
		zone := result.ZoneNumber
		band := Mgrs_latbands
}

func askBigQuery(w http.ResponseWriter, r *http.Request) {
	//firstValue, err := strconv.ParseFloat(r.FormValue("Latitude"), 64)
	//secondValue, err := strconv.ParseFloat(r.FormValue("Longtitute"), 64)
	//thirdValue:= r.FormValue("westLongditude")
	//fourthValue := r.FormValue("eastLongditude")
	// and use that context to create a new http client
	//client := http.DefaultClient
	// "bigquery-public-data:cloud_storage_geo_index.sentinel_2_index"
	//proj := os.Getenv("GOOGLE_CLOUD_PROJECT")
	ctx := context.Background()
	//fmt.Println(proj)
/*	tokenSource, err := google.DefaultTokenSource(Oauth2.NoContext, bigquery.BigqueryScope)
	if err != nil {
		log.Fatalf("Unable to acquire token source: %v", err)
	}
*/
 //scalabilitytest-183012

 	latLon := UTM.LatLon{50.77535, 6.08389}

	result, err := latLon.FromLatLon()
	if err != nil {
		panic(err.Error())
	}
	
		fmt.Printf(
			"Easting: %f; Northing: %f; ZoneNumber: %d; ZoneLetter: %s;",
			result.Easting,
			result.Northing,
			result.ZoneNumber,
			result.ZoneLetter,
		)
	
	test := "26.2199863395"
	client, err := bigquery.NewClient(ctx, "scalabilitytest-183012")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := fmt.Sprintf("SELECT granule_id, product_id FROM testgeoindex.sentinel_2_index_copy WHERE north_lat = %s LIMIT 1000",test)
	
	q := client.Query(query)
/*
	job, err := q.Run(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	jobID := job.ID()
	fmt.Printf("The job ID is %s\n", jobID)
	
	
	job, err = client.JobFromID(ctx, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Execute the query.
	it, err := job.Read(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Iterate through the results.
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(values)
	}
/*
	query := client.Query(
		`SELECT mgrs_tile
			FROM [bigquery-public-data:cloud_storage_geo_index.sentinel_2_index]
			WHERE north_lat = 26.2199863395 LIMIT 1000`)


/*
AND south_lat = ` + secondValue +`
			AND west_lon = ` + thirdValue +`
			AND east_lon = ` + fourthValue +`
*/

	// Execute the query.
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(" no go: failed at read query")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Iterate through the results.
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(" no go : failed at itr")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Granule_id: ", values[0])
		fmt.Println("Project_id: ", values[1])

		// print out project_id + granule_id trim off unwanted string, then make a call to GCS api:
		//"https://www.googleapis.com/storage/v1/b/gcp-public-data-sentinel-2/o/tiles%2F01%2FC%2FCV%2FS2A_MSIL1C_20160304T203515_N0201_R085_T01CCV_20160309T000729.SAFE%2FGRANULE%2FS2A_OPER_MSI_L1C_TL_SGS__20160305T043523_A003657_T01CCV_N02.01%2FIMG_DATA%2FS2A_OPER_MSI_L1C_TL_SGS__20160305T043523_A003657_T01CCV_B8A.jp2"
		// call function instead of println, that takes above params and makes a call to the GCS api just like below
	}
	/*
	//Lav url
	url := "http://storage.googleapis.com/"
	bucketName := "gcp-public-data-sentinel-2"
	objectName := "/tiles/01/C/CV/S2A_MSIL1C_20151221T205519_N0201_R028_T01CCV_20160329T181515.SAFE/GRANULE/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_N02.01/IMG_DATA/S2A_OPER_MSI_L1C_TL_EPA__20160325T184811_A002599_T01CCV_B02.jp2"
	ScopeDatastore := url + bucketName + objectName
				  "/tiles/12/R/VP/S2A_MSIL1C_20160207T180404_N0201_R141_T12RVP_20160208T025227.SAFE/GRANULE/S2A_OPER_MSI_L1C_TL_MPS__20160914T225212_A006430_T12RVP_N02.04/IMG_DATA/
	//lav query
	//q :=

	// now we can use that http client as before
	res, err := client.Get(ScopeDatastore)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Println(ioutil.ReadAll(res.Body))
	*/
}
