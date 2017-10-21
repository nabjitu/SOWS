package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	//"os"
	"math"
	"strings"
	"time"

	//"google.golang.org/appengine"
	//"google.golang.org/appengine/datastore"
	//"google.golang.org/appengine/urlfetch"
	"io/ioutil"

	// Imports the Google Cloud BigQuery client package.
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	//"cloud.google.com/go/storage"
	//"golang.org/x/oauth2/google"
	//"google.golang.org/api/compute/v1"
	"github.com/im7mortal/UTM"
	
)

//JSON struct, når vi får et JSON response vil vi gerne lave det om til et object vi kan manipulere og trække bestmte fields ud af, som vi kan returnere.
type JsonResponse struct {
	Kind  string `json:"kind"`
	Items []struct {
		Kind                    string    `json:"kind"`
		ID                      string    `json:"id"`
		SelfLink                string    `json:"selfLink"`
		Name                    string    `json:"name"`
		Bucket                  string    `json:"bucket"`
		Generation              string    `json:"generation"`
		Metageneration          string    `json:"metageneration"`
		ContentType             string    `json:"contentType"`
		TimeCreated             time.Time `json:"timeCreated"`
		Updated                 time.Time `json:"updated"`
		StorageClass            string    `json:"storageClass"`
		TimeStorageClassUpdated time.Time `json:"timeStorageClassUpdated"`
		Size                    string    `json:"size"`
		Md5Hash                 string    `json:"md5Hash"`
		MediaLink               string    `json:"mediaLink"`
		Crc32C                  string    `json:"crc32c"`
		Etag                    string    `json:"etag"`
	} `json:"items"`
}
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

func makeMGRS(lat float64, long float64) string {
	
//Function der laver LAT&LONG om til MGRS Koordinater

	Mgrs_e100k := [3]string{`ABCDEFGH`, `JKLMNPQR`, `STUVWXYZ`}

	Mgrs_n100k := [2]string{`ABCDEFGHJKLMNPQRSTUV`, `FGHJKLMNPQRSTUVABCDE`}

	latLon := UTM.LatLon{lat, long}
	
		result, err := latLon.FromLatLon()
		if err != nil {
			panic(err.Error())
		}
		
		
		zone := result.ZoneNumber)
		band := result.ZoneLetter

		col := int(math.Floor(result.Easting / 100000))
		e100k := Mgrs_e100k[(zone-1)%3]
		e100kcol := e100k[col-1]

		row := int(math.Floor(result.Northing / 100000))%20
		n100k := Mgrs_n100k[(zone-1)%2]
		n100krow := n100k[row]

		MGRS := fmt.Sprintf("%d %s %c%c",zone, band, e100kcol, n100krow)
		return MGRS
}


func askBigQuery(w http.ResponseWriter, r *http.Request) {
	firstValue, err := strconv.ParseFloat(r.FormValue("Latitude"),64)
	secondValue, err := strconv.ParseFloat(r.FormValue("Longtitute"),64)

	ctx := context.Background()
	defclient := http.DefaultClient
	//UTM Koordinater fra LAT&LONG koordinater
 	latLon := UTM.LatLon{firstValue, secondValue}
	// Lav MGRS koordinater fra LAT&LONG koordinater
	MGRS := makeMGRS(firstValue, secondValue) 

	//clean up string for use
	MGRSq := strings.Replace(MGRS," ", "", -1)

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
	//Opret forbindelse til Bigquery
	client, err := bigquery.NewClient(ctx, "scalabilitytest-183012")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Lav Bigquery query der finder info udfra MGRS
	query := fmt.Sprintf("SELECT base_url, granule_id, product_id FROM testgeoindex.sentinel_2_index_copy WHERE mgrs_tile = '%s' LIMIT 1000", MGRSq)
	
	q := client.Query(query)

	// Data structure til at håndtere Json fra Google API
	var jstruct []JsonResponse
	var resFinal []string

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
		url := fmt.Sprintf("%s",values[0])
		prefixbucket := strings.Replace(url, "gs://", "", 1)
		bucketsplit := strings.Split(prefixbucket,"/")
		bucket := fmt.Sprintf("%s",bucketsplit[0])
		resUrl := strings.Replace(prefixbucket,bucket+"/","",1)

		//Query til Google Storage API
		ScopeDatastore := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o?prefix=%s/GRANULE/%s/IMG_DATA",bucket,resUrl, values[1])

		//Kald til Google Storage API
		res, err := defclient.Get(ScopeDatastore)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Println("done")
		
		
		//Behandler API response fra Google Storage
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic (err.Error())
		}
		err := json.Unmarshal(body, &jstruct)
		if err != nil {
			fmt.Println("error:", err)
		}
		//fmt.Printf("json? : %s ",responsejson)
		resFinal = append(resFinal, jstruct.items.selfLink)



		//fmt.Println(ioutil.ReadAll(res.Body))
		
		
		//fmt.Println("Granule_id: ", values[0])
		//fmt.Println("Project_id: ", values[1])

		// print out project_id + granule_id trim off unwanted string, then make a call to GCS api:
		//"https://www.googleapis.com/storage/v1/b/gcp-public-data-sentinel-2/o/tiles%2F01%2FC%2FCV%2FS2A_MSIL1C_20160304T203515_N0201_R085_T01CCV_20160309T000729.SAFE%2FGRANULE%2FS2A_OPER_MSI_L1C_TL_SGS__20160305T043523_A003657_T01CCV_N02.01%2FIMG_DATA%2FS2A_OPER_MSI_L1C_TL_SGS__20160305T043523_A003657_T01CCV_B8A.jp2"
		// call function instead of println, that takes above params and makes a call to the GCS api just like below
	}
	fmt.Printf("json: %s ",resFinal)
	
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
