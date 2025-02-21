package events

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/aadarshvelu/bms/config"
)

var (
	producer sarama.SyncProducer
	topic    = "book_events"
)

type BookEvent struct {
	EventType string      `json:"event_type"` // for CREATE, UPDATE, DELETE
	BookID    uint        `json:"book_id"`
	Data      interface{} `json:"data"`
}

// InitKafka initializes the Kafka producer
func InitKafka() {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	brokers := []string{config.GetEnv("KAFKA_BROKER", "localhost:9092")}

	var err error
	producer, err = sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		log.Printf("Failed to initialize Kafka producer: %v", err)
	}
}

// PublishBookEvent sends a book event to Kafka
func PublishBookEvent(eventType string, bookID uint, data interface{}) error {
	event := BookEvent{
		EventType: eventType,
		BookID:    bookID,
		Data:      data,
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(payload),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	return nil
}
