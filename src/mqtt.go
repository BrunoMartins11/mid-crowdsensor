package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ReceivedData struct {
	DeviceID, Token  string
	ProbeData        []ProbeData
}

type ProbeData struct {
	MacAddress, Rssi string
	PrevDetected     int64 //in milliseconds
	Timestamp        *time.Time
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

func (data ReceivedData) addTimestampToProbes() {
	for i := range data.ProbeData {
		currentTime := time.Now()
		data.ProbeData[i].Timestamp = &currentTime
	}
}

func createMQTTClient() MQTT.Client {
	//create a ClientOptions struct setting the broker address, ClientId, turn
	//off trace output and set the default message handler
	opts := MQTT.NewClientOptions().AddBroker(os.Getenv("BROKER_URL"))
	opts.SetClientID(os.Getenv("CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_TOKEN"))
	opts.SetPassword("")
	opts.SetDefaultPublishHandler(f)
	client := MQTT.NewClient(opts)
	//create a client using the above ClientOptions
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
	return client
}

func subscribeTopic(client MQTT.Client) {
	//subscribe to the topic /go-mqtt/sample and request messages to be delivered
	//at a maximum qos of zero, wait for the receipt to confirm the subscription
	if token := client.Subscribe(os.Getenv("TOPIC"), 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for !flag {
		time.Sleep(1 * time.Second)
	}

	//unsubscribe from /go-mqtt/sample
	if token := client.Unsubscribe(os.Getenv("TOPIC")); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}

func publishToken(topic string, message string, client MQTT.Client) {
	token := client.Publish(topic, 0, false, message)

	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}