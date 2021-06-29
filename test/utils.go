package test

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"github.com/thanhpk/randstr"
	"log"
	"net/http"
	"os"
	"testing"
)
type Token struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

func SetupToken(t *testing.T) string{
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	response, err := http.Get(os.Getenv("IOT_AUTH")+ "/signup?username=" + randstr.String(16))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 406 {
		t.Fatal("Name of device exists")
	}
	token := new(Token)

	err = json.NewDecoder(response.Body).Decode(token)

	if err != nil {
		log.Fatal(err)
	}
	return token.Token
}

type MQClient struct {
	conn *amqp.Connection
}

func CreateMQClient() MQClient{
	return MQClient{}
}

func (mq MQClient) PublishToQueue(queueName string, payload []byte) {

	fmt.Println("Successfully Published Message to Queue" + queueName + string(payload))
}
