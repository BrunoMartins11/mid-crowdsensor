package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/api/option"
)

var flag bool = false

type ReceivedData struct {
	DeviceID  string
	ProbeData []ProbeData
}

type ProbeData struct {
	MacAddress, Rssi string
	PrevDetected     int64 //in miliseconds
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	var data ReceivedData
	err := json.Unmarshal(payload, &data)

	if err != nil {
		log.Fatalln(err)
	}

	if strings.Compare(string(payload), "\n") > 0 {
		fmt.Printf("TOPIC: %s\n", topic)
		fmt.Printf("MSG: %s\n", payload)
	}

	if strings.Compare("bye\n", string(payload)) == 0 {
		fmt.Println("exitting")
		flag = true
	}

	publishProbesToFirestore(data)

}

func connectToFirestore() (*firestore.Client, context.Context) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsFile("./crowdsensor-mid-storage-a90f7f257056.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client, ctx
}

func publishProbesToFirestore(probes ReceivedData) {

	client, ctx := connectToFirestore()
	for _, probe := range probes.ProbeData {

		_, _, err := client.Collection(probes.DeviceID).Add(ctx, probe)

		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
	}
}

func main() {
	//create a ClientOptions struct setting the broker address, clientid, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker("ssl://mqtt.flespi.io:8883")
	opts.SetClientID("Device-sub")
	opts.SetUsername("FlespiToken lCG8yJPUWRc9awe3M2AaTuKcqd5N4Nvgd1cByPklwkiuGmogcOgW6QWmURXOujSx")
	opts.SetPassword("")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe("cenas123", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for flag == false {
		time.Sleep(1 * time.Second)
	}

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("cenas123"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
