package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"os"
)

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

func publishProbesToFirestore(probes ReceivedData) {

	client, ctx := connectToFirestore()
	for _, probe := range probes.ProbeData {

		_, _, err := client.Collection(probes.DeviceID).Add(ctx, probe)

		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
	}
}
