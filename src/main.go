package main

import (
	"log"
	"net/http"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

var flag = false
var client MQTT.Client

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addDevice", addDeviceHandler)

	client = createMQTTClient()
	go subscribeTopic(client)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
