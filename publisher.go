package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// LogData represents the structure of the log message.
type LogData struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	IPAddress string `json:"ip_address"`
}

func main() {
	// NATS server URL (assuming it's running on localhost:4222)
	natsURL := "nats://localhost:4222" // Change if NATS server is running elsewhere

	// Connect to NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	fmt.Println("Connected to NATS")

	// Create an example log message
	logData := LogData{
		EventID:   uuid.New().String(), // Generate a unique event ID
		EventType: "login",
		UserID:    "12345",
		Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
		IPAddress: "192.168.1.1",
	}

	// Marshal the log data into JSON
	logJSON, err := json.Marshal(logData)
	if err != nil {
		log.Fatalf("Error marshaling log data: %v", err)
	}

	// Publish the log data to the NATS 'logs.audit' subject
	err = nc.Publish("logs.audit", logJSON)
	if err != nil {
		log.Fatalf("Error publishing to NATS: %v", err)
	}

	fmt.Println("Log data published to NATS")
}
