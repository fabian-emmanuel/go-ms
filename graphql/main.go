package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
)

type AppConfig struct {
	AccountUrl         string `envconfig:"ACCOUNT_SERVER_URL"`
	CatalogUrl         string `envconfig:"CATALOG_SERVER_URL"`
	OrderUrl           string `envconfig:"ORDER_SERVER_URL"`
	GraphQLServicePort int    `envconfig:"GRAPHQL_SERVICE_PORT"`
}

func main() {
	var config AppConfig
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewGraphQLServer(config.AccountUrl, config.CatalogUrl, config.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.New(s.ToExecutableSchema())
	srv.AddTransport(&transport.Websocket{})
	http.Handle("/graphql", srv)
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	log.Printf("Listening on port :%v...\n", config.GraphQLServicePort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", config.GraphQLServicePort), nil))
}
