package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/candrap89/droneTravelingDistance/repository"
	"github.com/segmentio/kafka-go"
)

// consumerHandler represents a handler for Kafka consumers.
type ConsumerHandler struct {
	Repository repository.RepositoryInterface
}

// define globa variables for this file
var (
	ackReceived = make(map[string]bool)
	ackMutex    = &sync.Mutex{}
)

func NewConsumerHandler(repository repository.RepositoryInterface) *ConsumerHandler {
	return &ConsumerHandler{Repository: repository}
}

func (ch *ConsumerHandler) StartNewProductConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "new-products-ack",
		GroupID:  "inventory-service",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var ack map[string]string
		json.Unmarshal(m.Value, &ack)
		fmt.Printf("Received message: %s\n", string(m.Value)) // {"cifId":"1234","status":"ok","timestamp":"2025-04-14T15:04:05Z"}

		userCif := ack["cifId"]

		fmt.Printf(" product: %+v\n", ack) //product map[cifId:1234 status:ok timestamp:2025-04-14T15:04:05Z]
		ackMutex.Lock()
		ackReceived[userCif] = true
		ackMutex.Unlock()
		fmt.Printf("Received estate: %s\n", ack["estate"])
		if ack["estate"] != "" {
			// Received estate: 1234
			// Create estate in the database
			// Assuming the ack contains length and width for the estate
			fmt.Printf("Received length: %s\n", ack["length"])
			fmt.Printf("Received width: %s\n", ack["width"])

			l, _ := strconv.Atoi(ack["length"])
			w, _ := strconv.Atoi(ack["width"])

			output, err := ch.Repository.CreateEstate(context.Background(), repository.CreateEstateInput{
				Width:  w,
				Length: l,
			})
			if err != nil {
				fmt.Printf("Failed to create estate: %v\n", err)
			}
			fmt.Printf("Estate created: %s\n", &output)
			r.Close()
		}

		fmt.Printf("Received ACK for product: %s\n", userCif) // Received ACK for product: 1234
	}

	r.Close()
}

func WaitForAck(UserCif string) bool {
	for {
		ackMutex.Lock()
		if ackReceived[UserCif] {
			delete(ackReceived, UserCif)
			ackMutex.Unlock()
			return true
		}
		ackMutex.Unlock()
	}
}
