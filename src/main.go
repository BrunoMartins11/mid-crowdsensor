package main

import (
	"fmt"
	//import the Paho Go MQTT library
	"os"
	"strings"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var flag bool = false

//var wcount int = 0

//define a function for the default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	topic := msg.Topic()
	payload := msg.Payload()
	if strings.Compare(string(payload), "\n") > 0 {
		fmt.Printf("TOPIC: %s\n", topic)
		fmt.Printf("MSG: %s\n", payload)
	}

	if strings.Compare("bye\n", string(payload)) == 0 {
		fmt.Println("exitting")
		flag = true
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
	if token := c.Subscribe("golangsub", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for flag == false {
		time.Sleep(1 * time.Second)
		//fmt.Println("waiting: ", wcount)
		//wcount += 1
	}

	//unsubscribe from /go-mqtt/sample
	if token := c.Unsubscribe("golangsub"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
}
