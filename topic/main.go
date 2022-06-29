package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

var (
	address, topic, auths string
	withAuth              bool
	partitions, replicas  int
)

func main() {
	flag.StringVar(&address, "kafkaURL", "kafka-headless:9092", "kafka url")
	flag.StringVar(&topic, "topic", "", "kafka topic")
	flag.StringVar(&auths, "auths", "", "kafka auths")
	flag.BoolVar(&withAuth, "withAuth", false, "kafka with auth")
	flag.IntVar(&partitions, "p", 3, "kafka topic partitions")
	flag.IntVar(&replicas, "r", 2, "kafka topic replicas")
	flag.Parse()

	if err := newKafkaTopic(address, topic, withAuth, auths, partitions, replicas); err != nil {
		log.Println(err)
	} else {
		log.Println(topic, "created!")
	}
}

func newKafkaTopic(address, topic string, withAuth bool, auths string, partitions, replicas int) error {
	// to create topics when auto.create.topics.enable='false'
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

	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	log.Println(net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	var controllerConn *kafka.Conn
	controllerConn, err = dialer.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicas,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}
