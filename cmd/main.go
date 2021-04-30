package main

import (
	"fmt"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/coms"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"log"
	"net/http"
	"os"
	"time"
)



func main() {
	fmt.Println(os.Getenv("PORT"))
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

	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
