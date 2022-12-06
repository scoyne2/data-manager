package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/scoyne2/data-manager-api/feed"
	"github.com/scoyne2/data-manager-api/schemas"
)

var API_HOST string = os.Getenv("API_HOST")
var FRONT_END_URL string = os.Getenv("FRONT_END_URL")

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
		fmt.Println(r.URL.String())
		fmt.Println(r.Method)
		fmt.Println(r.Header)
		frontEnd := fmt.Sprintf("http://%s", FRONT_END_URL)
		w.Header().Set("Access-Control-Allow-Origin", frontEnd)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		next.ServeHTTP(w,r)
    })
}

func check(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Health check</h1>")
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
	http.Handle("/graphql", CorsMiddleware(h))

	// Access via http://$API_HOST/sandbox
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
   initialEndpoint: "http://%s/graphql",
 });
 // advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
</html>`, API_HOST)

var sandboxHTML = []byte(sandboxHtmlString)

