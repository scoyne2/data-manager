package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/scoyne2/data-manager-api/feed"
	"github.com/scoyne2/data-manager-api/datapreview"
	"github.com/scoyne2/data-manager-api/schemas"
)

var API_HOST string = os.Getenv("API_HOST")
var FRONT_END_URL string = os.Getenv("FRONT_END_URL")
var DATA_PREVIEW_BUCKET string = os.Getenv("RESOURCES_BUCKET_NAME")

func main() {

	feedService := feed.NewService(
		feed.NewPostgressRepository(),
	)

	schema, err := schemas.GenerateSchema(&feedService)
	if err != nil {
		panic(err)
	}

	StartServer(schema)
}

func CorsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		frontEnd := fmt.Sprintf("https://%s", FRONT_END_URL)
		// localHost := "http://localhost.com:3000"
		w.Header().Set("Access-Control-Allow-Origin", frontEnd)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		next.ServeHTTP(w,r)
    })
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health Check</h1>")
}

// Example: https://api.datamanagertool.com/preview?vendor=Coyne+Enterprises&feedname=Hello&filename=hello.csv
func dataPreview(w http.ResponseWriter, r *http.Request){
	frontEnd := fmt.Sprintf("https://%s", FRONT_END_URL)
	w.Header().Set("Access-Control-Allow-Origin", frontEnd)
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	
	vals := r.URL.Query()
	vendor, ok := vals["vendor"]
	var vd string
	if ok {
		if len(vendor) >= 1 {
			vdLower := strings.ToLower(vendor[0])
			vd = strings.Replace(vdLower, " ", "_", -1)
		}
	}
	feedName, ok := vals["feedname"]
	var fn string
	if ok {
		if len(feedName) >= 1 {
			fnLower := strings.ToLower(feedName[0])
			fn = strings.Replace(fnLower, " ", "_", -1)
		}
	}
	fileName, ok := vals["filename"]
	var fln string
	if ok {
		if len(fileName) >= 1 {
			fln = fileName[0]
		}
	}

	results, err := datapreview.DataPreviewAthena(vd, fn, fln, DATA_PREVIEW_BUCKET)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	} 
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(results)
}

// StartServer will trigger the server with a Playground
func StartServer(schema *graphql.Schema) {
	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
		GraphiQL: false,
	})

	http.HandleFunc("/", check)
	http.HandleFunc("/health_check", check)
	http.HandleFunc("/preview", dataPreview)
	http.Handle("/graphql", CorsMiddleware(h))

	// Access via https://$API_HOST/sandbox
	http.Handle("/sandbox", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(sandboxHTML)
	}))

	fmt.Println("Server starting on port 8080...")
	fmt.Println(API_HOST)
	http.ListenAndServe(":8080", nil)
}


var sandboxHtmlString = fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
<div id="sandbox" style="height:100vh; width:100vw;"></div>
<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
<script>
 new window.EmbeddedSandbox({
   target: "#sandbox",
   // Pass through your server href if you are embedding on an endpoint.
   // Otherwise, you can pass whatever endpoint you want Sandbox to start up with here.
   initialEndpoint: "https://%s/graphql",
 });
 // advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
</html>`, API_HOST)

var sandboxHTML = []byte(sandboxHtmlString)

