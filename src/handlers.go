package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type AuthToken struct {
	TokenType string `json:"token_type"`
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
}

func addDeviceHandler(w http.ResponseWriter, req *http.Request) {
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
	token := new(AuthToken)
	json.NewDecoder(response.Body).Decode(token)

	if err != nil {
		log.Fatal(err)
	}

	go publishToken("BRUNO_ID4", token.Token, client)
}
