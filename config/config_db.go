package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"

	"github.com/KhoirulAziz99/final_project_e_wallet/pkg"
)

func InitDb() (*sql.DB, error) {
	host := pkg.GetEnv("DB_HOST")
	port := pkg.GetEnv("DB_PORT")
	user := pkg.GetEnv("DB_USER")
	password := pkg.GetEnv("DB_PASSWORD")
	dbname := pkg.GetEnv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	} else {
		log.Print("Succsessfully connected")
	}
	return db, nil
}

func DbClose(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	} else {
		log.Println("Databased closed")
	}
}
