package coms

import (
	"encoding/json"
	"fmt"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/auth"
	"github.com/BrunoMartins11/mid-crowdsensor/internal/status"
	"log"
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ReceivedData struct {
	DeviceID, Token string
	ProbeData       []status.ProbeData
}

type Fragment struct {
	Id, End int64 //End 0 if not last, 1 if last fragment
	Data    string
}

var flag = false

var Client MQTT.Client
var fragments = make(map[int64]string)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	var result map[string]interface{}

	_ = json.Unmarshal(payload, &result)
	if result["Id"] != nil {
		handleFragmentArrival(payload)
	} else if result["DeviceID"] != nil {
		createProbeData(payload, topic)
	}
}

func handleFragmentArrival(payload []byte) {
	var data Fragment

	err := json.Unmarshal(payload, &data)
	if err != nil {
		log.Fatalln(err)
	}

	frag, exists := fragments[data.Id]
	if exists && data.End == 0 {
		fragments[data.Id] = frag + data.Data
	} else if exists && data.End == 1 {
		createProbeDataFromFragments(frag + data.Data)
		delete(fragments, data.Id)
	} else {
		fragments[data.Id] = data.Data
	}
}

func createProbeData(payload []byte, topic string) {
	var data ReceivedData
	err := json.Unmarshal(payload, &data)
	if err != nil {
		log.Fatalln(err)
	}
	if auth.IsValidToken(PreToken + data.Token) {
		data.addTimestampToProbes()
		for _, val := range data.ProbeData{
			status.ManageNewProbe(val)
		}
		printTopicData(payload, topic)
		PublishProbesToFirestore(data)
	}
}

func createProbeDataFromFragments(payload string) {
	var data ReceivedData
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		log.Fatalln(err)
	}
	if auth.IsValidToken(PreToken + data.Token) {
		data.addTimestampToProbes()
		for _, val := range data.ProbeData{
			status.ManageNewProbe(val)
		}
		PublishProbesToFirestore(data)
	}
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
		data.ProbeData[i].Timestamp = currentTime
		data.ProbeData[i].DeviceID = data.DeviceID
	}
}

func CreateMQTTClient() MQTT.Client {
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

func SubscribeTopic(client MQTT.Client) {
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

func PublishToken(topic string, message string, client MQTT.Client) {
	token := client.Publish(topic, 0, false, message)

	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}
}
