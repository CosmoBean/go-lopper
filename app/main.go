package main

import (
	"go-lopper/db"
	"go-lopper/server"

	"github.com/joho/godotenv"
)

func main() {

	//init env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	//init db migrations
	db.Init()

	//init server
	server.SetupAndServe()
}
