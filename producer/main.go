package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var (
	kafkaURL string
	topic    string
	withAuth bool
	auths    string
	acks     int
)

func newKafkaWriter(kafkaURL, topic string, withAuth bool, auths string, acks int) *kafka.Writer {
	var sharedTransport *kafka.Transport
	sharedTransport = &kafka.Transport{
		DialTimeout: 10 * time.Second,
	}

	if withAuth {
		mechanism := plain.Mechanism{
			Username: strings.Split(auths, ":")[0],
			Password: strings.Split(auths, ":")[1],
		}
		sharedTransport = &kafka.Transport{
			DialTimeout: 10 * time.Second,
			SASL:        mechanism,
		}
	}

	kafkaURLs := strings.Split(kafkaURL, ",")
	return &kafka.Writer{
		Addr:         kafka.TCP(kafkaURLs...),
		Topic:        topic,
		RequiredAcks: kafka.RequiredAcks(acks),
		Balancer:     &kafka.RoundRobin{},
		Transport:    sharedTransport,
		BatchSize:    1000,
	}
}

func main() {
	flag.StringVar(&kafkaURL, "kafkaURL", "kafka-headless:9092", "kafka url")
	flag.StringVar(&topic, "topic", "", "kafka topic")
	flag.BoolVar(&withAuth, "withAuth", false, "kafka with auth")
	flag.StringVar(&auths, "auths", "admin:admin", "kafka auth")
	flag.IntVar(&acks, "acks", -1, "kafka acks")
	flag.Parse()

	writer := newKafkaWriter(kafkaURL, topic, withAuth, auths, acks)
	defer writer.Close()
	log.Println("start producing ... !!!")
	for i := 0; ; i++ {
		key := fmt.Sprintf("Key-%d", i)
		msg := kafka.Message{
			Key:   []byte(key),
			Value: []byte(fmt.Sprint(uuid.New())),
		}
		err := writer.WriteMessages(context.Background(), msg)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("produced", key)
		}
		time.Sleep(1 * time.Second)
	}
}
