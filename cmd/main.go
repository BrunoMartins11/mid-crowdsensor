package main

import (
	"github.com/BrunoMartins11/mid-crowdsensor/internal/coms"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)



func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addDevice", coms.AddDeviceHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	coms.Client = coms.CreateMQTTClient()
	status.InitializeRoomState()
	 go func() {
			for {
			time.Sleep(3*time.Minute)
			status.State.InitializeCleanup()
		}
	}()
	go coms.SubscribeTopic(coms.Client)

	log.Fatal(http.ListenAndServe(":1234", nil))
}
