package main

import (
	"github.com/fabian-emmanuel/go-ms/account"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseUrl        string `envconfig:"DATABASE_URL"`
	AccountServicePort int    `envconfig:"ACCOUNT_SERVICE_PORT"`
}

func main() {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	var repo account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		repo, err = account.NewPostgresRepository(config.DatabaseUrl)
		if err != nil {
			log.Println(err)
		}
		return
	})

	defer repo.Close()
	log.Printf("Listening on port :%v...\n", config.AccountServicePort)
	s := account.NewAccountService(repo)
	log.Fatal(account.ListenGRPC(s, config.AccountServicePort))
}
