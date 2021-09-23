package persistence

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

type ReceivedData struct {
	DeviceID, Token string
	ProbeData       []ProbeData
}

type ProbeData struct {
	DeviceID string
	MacAddress, Rssi string
	PrevDetected     int64 //in milliseconds
	Timestamp        time.Time
}

func connectToFirestore() (*firestore.Client, context.Context) {
	// Use a service account
	ctx := context.Background()
	sa := option.WithCredentialsJSON([]byte(os.Getenv("FIRESTORE_SECRETS")))
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

func PublishProbesToFirestore(probes []byte) {

	var data ReceivedData
	err := json.Unmarshal(probes, &data)
	if err != nil {
		log.Fatalln(err)
	}
	client, ctx := connectToFirestore()
	for _, probe := range data.ProbeData {

		_, _, err := client.Collection(probe.DeviceID).Add(ctx, probe)

		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
	}
}
