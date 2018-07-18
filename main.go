package main

import (
	"context"
	"flag"
	"log"

	"github.com/crossle/channel-father-mixin-bot/durable"
	"github.com/crossle/channel-father-mixin-bot/services"
)

func main() {
	service := flag.String("service", "http", "run a service")
	flag.Parse()
	db := durable.OpenDatabaseClient(context.Background())
	defer db.Close()

	switch *service {
	case "http":
		err := StartServer(db)
		if err != nil {
			log.Println(err)
		}
	default:
		hub := services.NewHub(db)
		err := hub.StartService(*service)
		if err != nil {
			log.Println(err)
		}
	}
}
