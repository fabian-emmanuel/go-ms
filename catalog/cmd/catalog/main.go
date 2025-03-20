package main

import (
	"github.com/fabian-emmanuel/go-ms/catalog"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseUrl        string `envconfig:"DATABASE_URL"`
	CatalogServicePort int    `envconfig:"CATALOG_SERVICE_PORT"`
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	var repo catalog.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		_, err = catalog.NewElasticRepository(config.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})

	log.Printf("Listening on port :%v...\n", config.CatalogServicePort)
	s := catalog.NewCatalogService(repo)
	log.Fatal(catalog.ListenGRPC(s, config.CatalogServicePort))
}
