// https://github.com/GoogleCloudPlatform/golang-samples/blob/master/appengine_flexible/pubsub/pubsub.go
package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/kolide/osquery-go"
	"github.com/kolide/osquery-go/plugin/logger"
	"log"
	"os"
)

var (
	c           context.Context
	client      *pubsub.Client
	socketPath  string
	topicName   string
	projectName string
)

func validateEnv(variable string) string {
	val, ok := os.LookupEnv(variable)

	if !ok {
		fmt.Printf("Error: Missing \"%s\" environment variable", variable)
		os.Exit(1)
	}

	return val
}

func main() {
	// Name of the GCP Project that contains the pub sub infra
	projectName = validateEnv("GCP_PROJECT")
	topicName = validateEnv("TOPIC")
	// Path to OSQUERY extensions socket ("/var/osquery/osquery.em")
	socketPath = validateEnv("SOCKET_PATH")
	// Path to JSON Google Cloud creds ("/var/osquery/certs/pubsub.json")
	validateEnv("GOOGLE_APPLICATION_CREDENTIALS")

	server, err := osquery.NewExtensionManagerServer("pubSubLogger", socketPath)
	if err != nil {
		log.Fatalf("Error creating extension: %s\n", err)
	}

	server.RegisterPlugin(logger.NewPlugin("pubSubLogger", LogString))
	err = server.Run()
	if err != nil {
		log.Println(err)
	}

}

func LogString(ctx context.Context, typ logger.LogType, logText string) error {
	// Dont log osquery status messages to Pub/Sub
	if typ == logger.LogTypeStatus {
		return nil
	}

	c = context.Background()
	client, err := pubsub.NewClient(c, projectName)

	if err != nil {
		log.Fatalln(err)
	}

	t := client.Topic(topicName)
	result := t.Publish(c, &pubsub.Message{Data: []byte(logText)})

	id, err := result.Get(c)
	if err != nil {
		return fmt.Errorf("Get: %v", err)
	}
	log.Printf("Published a message; msg ID: %v\n", id)
	return nil

	log.Printf("%s: %s\n", typ, logText)
	return nil
}
