package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	//"os"
	"math"
	"strings"
	"time"
	"bufio"

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
	"github.com/qedus/osmpbf"
	"github.com/golang/geo/s2"
)

//JSON struct, når vi får et JSON response vil vi gerne lave det om til et object vi kan manipulere og trække bestmte fields ud af, som vi kan returnere.
type JsonResponse struct {
	Kind  string `json:"kind"`
	Items []Item
}

type point struct {
	firstPoint float64
	secondPoint float64
}

type Item struct {
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
}

/*
type ReturnJson struct{
	ImageLinks				[]ImageLink
}
*/
type ImageLink struct {
	Link string `json:"link`
}

type StrResult struct{
	Strings [] string
}

func osmDecoder() {
	f, err := os.Open("greater-london-140324.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	var nc, wc, rc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				// Process Node v.
				nc++
			case *osmpbf.Way:
				// Process Way v.
				wc++
			case *osmpbf.Relation:
				// Process Relation v.
				rc++
			default:
				log.Fatalf("unknown type %T\n", v)
			}
		}
	}

	fmt.Printf("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)
}

func main() {

	/*p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(54.918, 8.552))
	p2 := s2.PointFromLatLng(s2.LatLngFromDegrees(55.048, 8.471))
	p3 := s2.PointFromLatLng(s2.LatLngFromDegrees(55.481, 12.736))
	p4 := s2.PointFromLatLng(s2.LatLngFromDegrees(54.837, 9.392))
	p5 := s2.PointFromLatLng(s2.LatLngFromDegrees(54.918, 8.552))

	points := []s2.Point{p1, p2, p3, p4, p5}

	fmt.Println("Points: ", points)

	l1 := s2.LoopFromPoints(points)

	fmt.Println("l1: ", l1)

	rect := l1.RectBound()

	fmt.Println("rect: ", rect)

	loops := []*s2.Loop{l1}

	fmt.Println("loops: ", loops[0])

	poly := s2.PolygonFromLoops(loops)

	fmt.Println("poly: ", poly.NumEdges())

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 100}
	cover := rc.Covering(poly)
	var c s2.Cell
	var totalArea float64
	totalArea = 0
	for i := 0; i < len(cover); i++ {
		fmt.Printf("Cell   %v   :   ", i)
		c = s2.CellFromCellID(cover[i])
		fmt.Printf("Low:   %v   -   ", c.RectBound().Lo())
		fmt.Printf("High:   %v   \n", c.RectBound().Hi())
		totalArea = totalArea + c.RectBound().Area()
	}

	fmt.Printf("Total   Area   with   multiple   rectangles:   %v", totalArea)*/

	//mux := http.NewServeMux()
	//files := http.FileServer(http.Dir(config.Static)) mux.Handle("/static/", http.StripPrefix("/static/", files))
	//mux.HandleFunc("/", index)

	http.HandleFunc("/", handler)
	http.HandleFunc("/a", askBigQuery)
	http.HandleFunc("/area", getArea)
	http.HandleFunc("/polyurl", readBodyFromUrl)
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

	latLon := UTM.LatLon{Latitude: lat, Longitude: long}

	result, err := latLon.FromLatLon()
	if err != nil {
		panic(err.Error())
	}

	zone := result.ZoneNumber
	band := result.ZoneLetter

	col := int(math.Floor(result.Easting / 100000))
	e100k := Mgrs_e100k[(zone-1)%3]
	e100kcol := e100k[col-1]

	row := int(math.Floor(result.Northing/100000)) % 20
	n100k := Mgrs_n100k[(zone-1)%2]
	n100krow := n100k[row]

	MGRS := fmt.Sprintf("%d %s %c%c", zone, band, e100kcol, n100krow)
	return MGRS
}

func askBigQuery(w http.ResponseWriter, r *http.Request) {
	firstValue, err := strconv.ParseFloat(r.FormValue("Latitude"), 64)
	secondValue, err := strconv.ParseFloat(r.FormValue("Longtitute"), 64)

	ctx := context.Background()
	defclient := http.DefaultClient
	//UTM Koordinater fra LAT&LONG koordinater
	latLon := UTM.LatLon{firstValue, secondValue}
	// Lav MGRS koordinater fra LAT&LONG koordinater
	MGRS := makeMGRS(firstValue, secondValue)

	//clean up string for use
	MGRSq := strings.Replace(MGRS, " ", "", -1)

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
	client, err := bigquery.NewClient(ctx, "bigquery-public-data")
	if err != nil {
		panic(err.Error())
	}
	//Lav Bigquery query der finder info udfra MGRS
	query := fmt.Sprintf("SELECT base_url, granule_id, product_id FROM cloud_storage_geo_index.sentinel_2_index WHERE mgrs_tile = '%s'", MGRSq)

	q := client.Query(query)

	// Data structure til at håndtere Json fra Google API

	var resFinal []byte
	//resFinal := make([]string, "")

	// Execute the query.
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(" no go: failed at read query")
		panic(err.Error())
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
			panic(err.Error())
		}
		url := fmt.Sprintf("%s", values[0])
		prefixbucket := strings.Replace(url, "gs://", "", 1)
		bucketsplit := strings.Split(prefixbucket, "/")
		bucket := fmt.Sprintf("%s", bucketsplit[0])
		resUrl := strings.Replace(prefixbucket, bucket+"/", "", 1)

		//Query til Google Storage API
		ScopeDatastore := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o?prefix=%s/GRANULE/%s/IMG_DATA", bucket, resUrl, values[1])

		//Kald til Google Storage API
		res, err := defclient.Get(ScopeDatastore)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println("done")

		//Behandler API response fra Google Storage
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}

		var jstruct JsonResponse
		json.Unmarshal(body, &jstruct)
		length := len(jstruct.Items)

		for i := 0; i < length; i++ {
			t := &ImageLink{Link: jstruct.Items[i].SelfLink}
			b, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			//fmt.Println(string(b))
			resFinal = append(resFinal, b...)
		}

	}

	//fmt.Printf("json: %s",resFinal)
	fmt.Fprintf(w, "Json: %s", resFinal)
}

func getArea(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	firstValue, err := strconv.ParseFloat(r.FormValue("Latitude"), 64)
	secondValue, err := strconv.ParseFloat(r.FormValue("Longtitute"), 64)
	thirdValue, err := strconv.ParseFloat(r.FormValue("Latitude2"), 64)
	fourthValue, err := strconv.ParseFloat(r.FormValue("Longtitute2"), 64)
	if err != nil {
		// handle error
	}

	//MGRS1 := makeMGRS(firstValue, secondValue)

	//MGRS2 := makeMGRS(thirdValue, fourthValue)

	//fmt.Println(MGRS1)
	//fmt.Println(MGRS2)

	//regn 3 og 4 punkt ud
	//coord 3 = firstvalue + fourth value
	//coord 4 = thirdvalue + second value
	//var Mgrsarea []string
	Mgrsarea := NewStringSet()

	if firstValue < thirdValue {
		for i := firstValue; i <= thirdValue; i += 0.1 {
			fmt.Println(i, secondValue)
			MGRS := makeMGRS(i, secondValue)

			MGRSs := strings.Replace(MGRS, " ", "", -1)

			Mgrsarea.Add(MGRSs)

			if secondValue < fourthValue {
				for j := secondValue; j <= fourthValue; j += 0.1 {
					MGRSj := makeMGRS(i, j)
					MGRSsj := strings.Replace(MGRSj, " ", "", -1)

					Mgrsarea.Add(MGRSsj)

				}
			} else if fourthValue < secondValue {
				for j := secondValue; j >= fourthValue; j -= 0.1 {
					MGRSj := makeMGRS(i, j)
					MGRSsj := strings.Replace(MGRSj, " ", "", -1)

					Mgrsarea.Add(MGRSsj)

				}
			}
		}
	} else if thirdValue < firstValue {
		for i := thirdValue; i <= firstValue; i += 0.1 {

			MGRS := makeMGRS(i, fourthValue)

			MGRSs := strings.Replace(MGRS, " ", "", -1)

			Mgrsarea.Add(MGRSs)

			if fourthValue < secondValue {
				for j := fourthValue; j <= secondValue; j += 0.1 {
					MGRSj := makeMGRS(i, j)
					MGRSsj := strings.Replace(MGRSj, " ", "", -1)

					Mgrsarea.Add(MGRSsj)

				}
			} else if secondValue < fourthValue {
				for j := fourthValue; j >= secondValue; j -= 0.1 {
					MGRSj := makeMGRS(i, j)
					MGRSsj := strings.Replace(MGRSj, " ", "", -1)

					Mgrsarea.Add(MGRSsj)

				}
			}
		}

	}

	//fmt.Println(Mgrsarea)

	c := make(chan string)
	Mgrsareaslice := Mgrsarea.AsSlice()
	for i := 0; i < len(Mgrsareaslice); i++ {
		//fmt.Println(Mgrsareaslice[i])
		go hellos(Mgrsareaslice[i], c)
		s := <-c
		fmt.Fprintf(w, s)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("done")

	fmt.Println("///////////////BENCHMARK END///////////////////")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("QUery time was ", elapsed)
	fmt.Println("///////////////BENCHMARK END///////////////////")

}

func hellos(s string, c chan string) {

	ctx := context.Background()
	defclient := http.DefaultClient

	//Opret forbindelse til Bigquery
	client, err := bigquery.NewClient(ctx, "sill-179308")
	if err != nil {
		panic(err.Error())
	}
	//Lav Bigquery query der finder info udfra MGRS

	query := fmt.Sprintf("SELECT base_url, granule_id, product_id FROM dataset_s.sentinel_2_index_copy WHERE mgrs_tile = '%s'", s)

	q := client.Query(query)

	// Data structure til at håndtere Json fra Google API

	var resFinal []byte
	//resFinal := make([]string, "")

	// Execute the query.
	it, err := q.Read(ctx)
	if err != nil {
		fmt.Println(" no go: failed at read query")
		panic(err.Error())
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
			panic(err.Error())
		}
		url := fmt.Sprintf("%s", values[0])
		prefixbucket := strings.Replace(url, "gs://", "", 1)
		bucketsplit := strings.Split(prefixbucket, "/")
		bucket := fmt.Sprintf("%s", bucketsplit[0])
		resUrl := strings.Replace(prefixbucket, bucket+"/", "", 1)

		//Query til Google Storage API
		ScopeDatastore := fmt.Sprintf("https://www.googleapis.com/storage/v1/b/%s/o?prefix=%s/GRANULE/%s/IMG_DATA", bucket, resUrl, values[1])

		//Kald til Google Storage API
		res, err := defclient.Get(ScopeDatastore)
		if err != nil {
			panic(err.Error())
		}
		//fmt.Println("done")

		//Behandler API response fra Google Storage
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}

		var jstruct JsonResponse
		json.Unmarshal(body, &jstruct)
		length := len(jstruct.Items)

		for i := 0; i < length; i++ {
			t := &ImageLink{Link: jstruct.Items[i].SelfLink}
			b, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			//fmt.Println(string(b))
			resFinal = append(resFinal, b...)
		}

	}
	c <- string(resFinal)
	//fmt.Println(string(resFinal))

}

type StringSet map[string]bool

func NewStringSet() StringSet {
	return make(StringSet)
}
func (s StringSet) Add(val string) {
	s[val] = true
}
func (s StringSet) AddAll(src StringSet) {
	for k, _ := range src {
		s[k] = true
	}
}
func (s StringSet) String() string {
	return fmt.Sprint(s.AsSlice()) // could be made more efficient if needed
}

func (s StringSet) AsSlice() []string {
	ret := make([]string, 0, len(s))
	for k, _ := range s {
		ret = append(ret, k)
	}
	return ret
}

func readBodyFromUrl(w http.ResponseWriter, r *http.Request) {
	Continent := strings.ToLower(r.FormValue("Continent"))
	Country := strings.ToLower(r.FormValue("Country"))
	Region := strings.ToLower(r.FormValue("Region"))

	//Tag Body fra URL'en

	url := fmt.Sprintf("")
	/*
		if Region == "" {
			req, err := http.NewRequest("GET", url2, nil)
			if err != nil {
				log.Fatalf("could not create request: %v", err)
			}

		} else if Country == "" && Region == "" {
			req, err := http.NewRequest("GET", url3, nil)
			if err != nil {
				log.Fatalf("could not create request: %v", err)
			}
		} else {
			req, err := http.NewRequest("GET", url1, nil)
			if err != nil {
				log.Fatalf("could not create request: %v", err)
			}
		}
	*/

	if Region == "" && Country != "" {
		url = fmt.Sprintf("http://download.geofabrik.de/%s/%s.poly", Continent, Country)

	} else if Region == "" && Country == "" {
		url = fmt.Sprintf("http://download.geofabrik.de/%s.poly", Continent)

	} else {
		url = fmt.Sprintf("http://download.geofabrik.de/%s/%s/%s.poly", Continent, Country, Region)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("could not create request: %v", err)
	}
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("http request failed: %v", err)
	}

	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var points []float64
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	i := 0
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		if i > 2{
			l := scanner.Text()
			
			if l == "END"{
				i++
			}else {
				ll, err := strconv.ParseFloat(l, 64)
				if err != nil {
					panic(err.Error())
				}
				points = append(points, ll)
				i++
				
			}
			
		}else {			
			i++
		}
	}

	var allPoints []s2.Point

	for i := 0 ; i <= len(points)/2 ; i= i+2 {
		p1 := s2.PointFromLatLng(s2.LatLngFromDegrees(points[i], points[i+1]))
		allPoints = append(allPoints, p1)
		fmt.Println(p1)
	}

	l1 := s2.LoopFromPoints(allPoints)
	rect := l1.RectBound()
	loops := []*s2.Loop{l1}
	poly := s2.PolygonFromLoops(loops)
	fmt.Printf("No. of edges %v\n", poly.NumEdges())

	//   one   big   rectangle   bounding   box,   just   to   test
	rect = poly.RectBound()
	fmt.Printf("Rect. Lat. Lo: %v \n", rect.Lat.Lo*180.0/math.Pi)
	fmt.Printf("Rect. Lat. Hi: %v \n", rect.Lat.Hi*180.0/math.Pi)
	fmt.Printf("Rect. Lng. Lo: %v \n", rect.Lng.Lo*180.0/math.Pi)
	fmt.Printf("Rect. Lng. Hi: %v \n", rect.Lat.Hi*180.0/math.Pi)
	fmt.Printf("\nOne Big Rect. Area %v\n\n", rect.Area())

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 100}
	cover := rc.Covering(poly)
	var c s2.Cell
	var totalArea float64
	totalArea = 0
	for i := 0; i < len(cover); i++ {
		fmt.Printf("Cell %v : ", i)
		c = s2.CellFromCellID(cover[i])
		fmt.Printf("Low: %v - ", c.RectBound().Lo())
		fmt.Printf("High: %v \n", c.RectBound().Hi())
		totalArea = totalArea + c.RectBound().Area()
	}
	fmt.Printf("Total Area with multiple rectangles: %v", totalArea)
	
	

	//fmt.Println(points)
	//fmt.Println(url)
	fmt.Println(res.Status)
	//fmt.Println(string(body))
	//fmt.Printf("%s", temp)
	res.Body.Close()
}
