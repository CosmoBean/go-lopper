package main

import (
	"go-lopper/model"

	"github.com/joho/godotenv"
)

func main() {

	//init env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	//init db migrations
	model.Setup()
}
