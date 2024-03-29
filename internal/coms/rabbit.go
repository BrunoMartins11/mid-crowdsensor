package coms

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type MQClient struct {
	conn *amqp.Connection
}

func CreateMQClient() MQClient{
	conn, err := amqp.Dial(os.Getenv("RABBIT_URL"))

	if err != nil {
		log.Fatal(err)
	}
	return MQClient{conn}
}

func (mq MQClient) PublishToQueue(queueName string, payload []byte) {
	ch, err := mq.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: payload,
		},
	)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}