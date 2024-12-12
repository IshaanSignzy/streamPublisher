package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// LogData represents the structure of the log message.
type LogData struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
	IPAddress string `json:"ip_address"`
}

// Function to simulate publishing log data to NATS
func publishToNATS(nc *nats.Conn, wg *sync.WaitGroup, totalMessages int) {
	defer wg.Done()

	// Loop to publish the given number of messages
	for i := 0; i < totalMessages; i++ {
		// Generate random log data
		logData := LogData{
			EventType: randomEventType(),
			UserID:    fmt.Sprintf("user-%d", rand.Intn(1000)), // Random user ID
			Timestamp: time.Now().UTC().Format("2006-01-02 15:04:05"),
			IPAddress: fmt.Sprintf("192.168.1.%d", rand.Intn(255)), // Random IP address
		}

		// Marshal the log data into JSON
		logJSON, err := json.Marshal(logData)
		if err != nil {
			log.Printf("Error marshaling log data: %v", err)
			return
		}

		// Publish the log data to the NATS 'logs.audit' subject
		err = nc.Publish("logs.audit", logJSON)
		if err != nil {
			log.Printf("Error publishing to NATS: %v", err)
			return
		}

		// Print a log every 100 messages to see progress
		if i%100 == 0 {
			fmt.Printf("Published %d messages\n", i)
		}

		// Optional: Introduce a delay to control the message rate
		// You can adjust the sleep duration to simulate different message frequencies
		time.Sleep(time.Millisecond * 10) // Simulate sending messages at a rate of 100 per second
	}
}

// Randomly generate an event type for testing
func randomEventType() string {
	events := []string{"login", "logout", "purchase", "view", "signup"}
	return events[rand.Intn(len(events))]
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

	// Total number of messages to publish in the load test
	totalMessages := 10000     // Adjust for load test volume
	concurrentPublishers := 10 // Number of concurrent Goroutines

	// WaitGroup to wait for all Goroutines to complete
	var wg sync.WaitGroup

	// Start multiple Goroutines to simulate concurrent publishing
	for i := 0; i < concurrentPublishers; i++ {
		wg.Add(1)
		go publishToNATS(nc, &wg, totalMessages/concurrentPublishers)
	}

	// Wait for all Goroutines to finish
	wg.Wait()

	fmt.Println("Load test completed: All messages have been published.")
}
