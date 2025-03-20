package main

import (
	"github.com/fabian-emmanuel/go-ms/order"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseUrl       string `envconfig:"DATABASE_URL"`
	AccountServiceUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogServiceUrl string `envconfig:"CATALOG_SERVICE_URL"`
	OrderServicePort  int    `envconfig:"ORDER_SERVICE_PORT"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	var repo order.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repo, err = order.NewPostgresRepository(config.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})

	defer repo.Close()
	log.Printf("Listening on port :%v...\n", config.OrderServicePort)
	s := order.NewOrderService(repo)
	log.Fatal(order.ListenGRPC(s, config.AccountServiceUrl, config.CatalogServiceUrl, config.OrderServicePort))
}
