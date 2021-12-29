package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
)

func main() {
	log.Println("Starting...")

	name := os.Getenv("INSTANCE_NAME")
	if name == "" {
		log.Fatal("The environment variable INSTANCE_NAME is not defined\nExitting")
	}

	log.Printf("Instance name: %s\n", name)

	projectID := os.Getenv("PROJECT_ID")
	if projectID == "" {
		log.Fatal("The environment variable PROJECT_ID is not defined\nExitting")
	}

	log.Printf("Project ID: %s\n", projectID)

	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	log.Println("GCP Firestore client created")

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)
		log.Println("Listening for shutdown signals")

		<-sigint
		log.Println("Shudown signal received")

		if err := client.Close(); err != nil {
			log.Fatalf("Failed to close client: %v", err)
		}

		log.Println("Client closed")
		log.Println("Exitting")

		os.Exit(0)
	}()

	for {
		now := time.Now()
		log.Printf("Got time: %v\n", now)

		_, err = client.Collection("uptime").Doc(name).Set(ctx, map[string]interface{}{
			"last_up": now,
		}, firestore.MergeAll)

		log.Println("Updated database")

		if err != nil {
			log.Fatalf("Failed to update document: %v", err)
		}

		time.Sleep(5 * time.Minute)
	}

}
