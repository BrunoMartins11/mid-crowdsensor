package main

import (
	"fmt"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/api"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/coms"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	http.HandleFunc("/addDevice", api.AddDeviceHandler)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

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

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
