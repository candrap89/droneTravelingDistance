package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

func SendNewProductMessage(barcode string) {
	conn, _ := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "new-products", 0)
	defer conn.Close()

	message := map[string]string{
		"barcode": barcode,
	}

	msg, _ := json.Marshal(message)
	conn.WriteMessages(
		kafka.Message{Value: msg},
	)
}
