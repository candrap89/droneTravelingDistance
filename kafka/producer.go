package kafka

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"
)

func SendNewProductMessage(votes int, cityName string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "city-votes", 0)
	if err != nil {
		return err
	}
	defer conn.Close()

	message := map[string]string{
		"votes":    strconv.Itoa(votes),
		"cityName": cityName,
	}

	msg, _ := json.Marshal(message)
	conn.WriteMessages(
		kafka.Message{Value: msg},
	)
	return err
}
