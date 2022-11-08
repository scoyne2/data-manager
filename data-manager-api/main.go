package main

import (
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/scoyne2/data-manager-api/feed"
	"github.com/scoyne2/data-manager-api/schemas"

)

func main() {

	feedService := feed.NewService(
		// TODO connect to postgress DB
		feed.NewMemoryRepository(),
	)

	schema, err := schemas.GenerateSchema(&feedService)
	if err != nil {
		panic(err)
	}

	StartServer(schema)
}

// StartServer will trigger the server with a Playground
func StartServer(schema *graphql.Schema) {
	h := handler.New(&handler.Config{
		Schema: schema,
		Pretty: true,
		GraphiQL: false,
	})

	http.Handle("/graphql", h)

	// Access via http://localhost:8080/sandbox
	http.Handle("/sandbox", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(sandboxHTML)
	}))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

var sandboxHTML = []byte(`
<!DOCTYPE html>
<html lang="en">
<body style="margin: 0; overflow-x: hidden; overflow-y: hidden">
<div id="sandbox" style="height:100vh; width:100vw;"></div>
<script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
<script>
 new window.EmbeddedSandbox({
   target: "#sandbox",
   // Pass through your server href if you are embedding on an endpoint.
   // Otherwise, you can pass whatever endpoint you want Sandbox to start up with here.
   initialEndpoint: "http://localhost:8080/graphql",
 });
 // advanced options: https://www.apollographql.com/docs/studio/explorer/sandbox#embedding-sandbox
</script>
</body>
</html>`)