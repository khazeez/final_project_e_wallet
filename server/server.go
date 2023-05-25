package server

import (
	"log"

	"github.com/KhoirulAziz99/final_project_e_wallet/api"
	"github.com/KhoirulAziz99/final_project_e_wallet/config"
	"github.com/KhoirulAziz99/final_project_e_wallet/pkg"
)

func Run() error {
	db, err := config.InitDb()
	if err != nil {
		panic(err)
	}

	defer config.DbClose(db)

	router := api.SetUpRouter(db)

	serverAddress := pkg.GetEnv("SERVER_ADDRESS")
	log.Printf("Server is running on address : %s \n", serverAddress)

	if err := router.Run(serverAddress); err != nil {
		return err
	}

	return nil
}