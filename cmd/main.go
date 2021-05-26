package main

import (
	"fmt"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/coms"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

type MSG struct {
	DeviceID string
	MacAddress string
	Active     bool //in milliseconds
	Timestamp        time.Time
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addDevice", coms.AddDeviceHandler)

	status.InitializeRoomState(coms.CreateMQClient())

	coms.Client = coms.CreateMQTTClient()

	 go func() {
	 	for {
			time.Sleep(3*time.Minute)
			status.State.InitializeCleanup()
			fmt.Println("Cleanup")
		}
	}()
	go coms.SubscribeTopic(coms.Client)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
