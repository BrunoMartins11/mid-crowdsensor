package main

import (
	"github.com/BrunoMartins11/mid-crowdsensor/src/coms"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addDevice", coms.AddDeviceHandler)

	coms.Client = coms.CreateMQTTClient()
	go coms.SubscribeTopic(coms.Client)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
