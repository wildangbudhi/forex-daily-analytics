package main

import (
	"log"

	"github.com/wildangbudhi/forex-daily-analytics/utils"
)

func main() {
	server, err := utils.NewServer()

	if err != nil {
		log.Fatal(err)
	}

	defer server.DB.Disconnect(*server.DBContext)

	server.Router.Run()
}
