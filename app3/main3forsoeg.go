package app3
//Natalie leger

import (
//"encoding/json"
"fmt"
"net/http"
"strconv"
//"os"
"math"
"strings"

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



func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/a", askBigQuery)
	//http.HandleFunc("/jsons", decodeHandler)
	http.ListenAndServe(":9000", nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
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

func Mgrs struct (zone, band, e100k, n100k, easting, northing, datum)(mgrs Mgrs) {datum

if (datum === undefined) datum = LatLon.datum.WGS84; // default if not supplied

if (!(1<=zone && zone<=60)) throw new Error('Invalid MGRS grid reference (zone ‘'+zone+'’)');
if (band.length != 1) throw new Error('Invalid MGRS grid reference (band ‘'+band+'’)');
if (Mgrs.latBands.indexOf(band) == -1) throw new Error('Invalid MGRS grid reference (band ‘'+band+'’)');
if (e100k.length!=1) throw new Error('Invalid MGRS grid reference (e100k ‘'+e100k+'’)');
if (n100k.length!=1) throw new Error('Invalid MGRS grid reference (n100k ‘'+n100k+'’)');

zone: = Number(zone);
band := band;
e100k := e100k;
n100k := n100k;
easting := Number(easting);
northing := Number(northing);
datum := datum;
}


func makeMGRS(lat float64, long float64) string {


	fmt.Printf("lat : %f",lat)
	fmt.Printf("long : %f", long)

	//Mgrs_latbands := `CDEFGHIJKLMNPQRSTUVWXX`
	//latbands := []rune(Mgrs_latbands)

	Mgrs_e100k := [3]string{`ABCDEFGH`, `JKLMNPQR`, `STUVWXYZ`}
	//e100k := []rune(Mgrs_e100k)

	Mgrs_n100k := [2]string{`ABCDEFGHJKLMNPQRSTUV`, `FGHJKLMNPQRSTUVABCDE`}
	//n100k := []rune(Mgrs_n100k)

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
	//fmt.Printf("lat: %f",lat)
	//fmt.Printf("math: %f",math.Floor(lat/8.0+10.0))
	band := result.ZoneLetter//latbands[int(math.Floor(lat/8.0+10.0))]

	col := int(math.Floor(result.Easting / 100000))
	e100k := Mgrs_e100k[(zone-1)%3]
	e100kcol := e100k[col-1]

	row := int(math.Floor(result.Northing / 100000))%20
	n100k := Mgrs_n100k[(zone-1)%2]
	n100krow := n100k[row]

	//fmt.Printf("%d %s %c%c",zone, band, e100kcol, n100krow)
	MGRS := fmt.Sprintf("%d %s %c%c",zone, band, e100kcol, n100krow)
	return MGRS
}

func askBigQuery(w http.ResponseWriter, r *http.Request) {
	firstValue, err := strconv.ParseFloat(r.FormValue("Latitude"),64)
	secondValue, err := strconv.ParseFloat(r.FormValue("Longtitute"),64)
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


	//fmt.Printf("latfirst: %f",firstValue)
	latLon := UTM.LatLon{firstValue, secondValue}

	MGRS := makeMGRS(firstValue, secondValue)
	fmt.Printf("MGRS: %s ",MGRS)
	MGRSq := strings.Replace(MGRS," ", "", -1)

	MGRSs := strings.Split(MGRS, " ")
	fmt.Printf("MGRSq: = %s",MGRSq)

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


	/*
		ScopeDatastore := fmt.Sprintf("http://legallandconverter.com/cgi-bin/android5c.cgi?username=DEVELOPX&password=TEST1234&latitude=%f&longitude=%f&cmd=mgrsrev1",firstValue,secondValue)


			res, err := http.Get(ScopeDatastore)
			if err != nil {
				http.Error(w, fmt.Sprintf("could not get google: %v", err), http.StatusInternalServerError)
				return
			}
			fmt.Println("MGRS: ", ioutil.ReadAll(res.Body))
	*/
	//test := "26.2199863395"
	client, err := bigquery.NewClient(ctx, "scalabilitytest-183012")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	query := fmt.Sprintf("SELECT base_url, granule_id, product_id FROM testgeoindex.sentinel_2_index_copy WHERE mgrs_tile = '%s' LIMIT 1000", MGRSq)
	//query := fmt.Sprintf("SELECT granule_id, product_id FROM testgeoindex.sentinel_2_index_copy WHERE north_lat = %s LIMIT 1000",test)
	q := client.Query(query)

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
		fmt.Printf("http://storage.googleapis.com/gcp-public-data-sentinel-2/tiles/%s/%s/%s/%s/GRANULE/%s/IMG_DATA/",MGRSs[0],MGRSs[1], MGRSs[2], values[1], values[2])


	}
}
