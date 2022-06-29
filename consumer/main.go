package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var (
	kafkaURL string
	topic    string
	withAuth bool
	auths    string
	groupID  string
)

func getKafkaReader(kafkaURL, topic string, withAuth bool, auths, groupID string) *kafka.Reader {
	var dialer *kafka.Dialer
	dialer = &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	if withAuth {
		mechanism := plain.Mechanism{
			Username: strings.Split(auths, ":")[0],
			Password: strings.Split(auths, ":")[1],
		}
		dialer = &kafka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
		}
	}

	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		Dialer:   dialer,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func main() {
	flag.StringVar(&kafkaURL, "kafkaURL", "kafka-headless:9092", "kafka url")
	flag.StringVar(&topic, "topic", "", "kafka topic")
	flag.BoolVar(&withAuth, "withAuth", false, "kafka with auth")
	flag.StringVar(&auths, "auths", "admin:admin", "kafka auth")
	flag.StringVar(&groupID, "groupID", "cgi", "kafka consumer-group-id")
	flag.Parse()

	reader := getKafkaReader(kafkaURL, topic, withAuth, auths, groupID)

	defer reader.Close()

	log.Println("start consuming ... !!!")
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}
