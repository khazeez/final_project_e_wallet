package main

import (
	"log"

	"github.com/KhoirulAziz99/final_project_e_wallet/server"
)

func main() {

	if err := server.Run; err != nil {
		log.Fatalf("Error running the server : %s", err())
	}
}
