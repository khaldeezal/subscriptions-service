package main

import (
	"log"

	"github.com/khaldeezal/subscriptions-service/cmd/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
