package main

import (
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strings"
	"time"
)

type ReceivedData struct {
	DeviceID  string
	ProbeData []ProbeData
}

type ProbeData struct {
	MacAddress, Rssi string
	PrevDetected     int64 //in milliseconds
	Timestamp *time.Time
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	var data ReceivedData
	err := json.Unmarshal(payload, &data)

	data.addTimestampToProbes()

	if err != nil {
		log.Fatalln(err)
	}

	printTopicData(payload, topic)

	publishProbesToFirestore(data)
}

func printTopicData(payload []byte, topic string) {
	if strings.Compare(string(payload), "\n") > 0 {
		fmt.Printf("TOPIC: %s\n", topic)
		fmt.Printf("MSG: %s\n", payload)
	}

	if strings.Compare("bye\n", string(payload)) == 0 {
		fmt.Println("exiting")
		flag = true
	}
}

func (data ReceivedData)addTimestampToProbes() {
	for i := range data.ProbeData {
		currentTime := time.Now()
		data.ProbeData[i].Timestamp = &currentTime
	}
}

func createMQTTClient (){
	//create a ClientOptions struct setting the broker address, ClientId, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(os.Getenv("BROKER_URL"))
	opts.SetClientID(os.Getenv("CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_TOKEN"))
	opts.SetPassword("")
	opts.SetDefaultPublishHandler(f)

	//create and start a client using the above ClientOptions
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := c.Subscribe(os.Getenv("TOPIC"), 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for !flag {
		time.Sleep(1 * time.Second)
	}

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe(os.Getenv("TOPIC")); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
