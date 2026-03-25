package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"

	"crm/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

//go:embed frontend.html
var frontendFS embed.FS

func main() {
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))

	http.Handle("/query", srv)
	http.HandleFunc("/", playground.Handler("CRM GraphQL", "/query"))

	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data, _ := frontendFS.ReadFile("frontend.html")
		w.Write(data)
	})

	fmt.Println("GraphQL server running on http://localhost:8080")
	fmt.Println("Playground: http://localhost:8080/")
	fmt.Println("Frontend: http://localhost:8080/app")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
