package main

import (
	"github.com/joho/godotenv"
	"log"
)

var flag = false


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	createMQTTClient()
}
