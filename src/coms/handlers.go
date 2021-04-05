package coms

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

type Token struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

var PreToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9."


func AddDeviceHandler(w http.ResponseWriter, req *http.Request) {
	device := req.URL.Query().Get("device")

	response, err := http.Get(os.Getenv("IOT_AUTH") + "/signup?username=" + device)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode == 406 {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	token := new(Token)

	err = json.NewDecoder(response.Body).Decode(token)

	if err != nil {
		log.Fatal(err)
	}
	msg := strings.Split(token.Token, ".")

	go PublishToken(device, msg[1]+"."+msg[2], Client)

	w.WriteHeader(http.StatusOK)
}
