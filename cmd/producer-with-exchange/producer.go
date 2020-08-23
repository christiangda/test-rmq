package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var (
	serverAddress string
	serverPort    int
	username      string
	password      string

	virtualHost        string
	exchange           string
	exchangeRoutingKey string
	exchangeType       string
	exchangeDurable    bool
	exchangeAutoDelete bool

	message            string
	messageContentType string
)

func main() {
	flag.StringVar(&serverAddress, "server-address", "127.0.0.1", "RabbitMQ server IP address")
	flag.IntVar(&serverPort, "server-port", 5672, "RabbitMQ server port")
	flag.StringVar(&username, "username", "guest", "RabbitMQ username")
	flag.StringVar(&password, "password", "guest", "RabbitMQ password")

	flag.StringVar(&virtualHost, "virtualHost", "my.virtualhost", "RabbitMQ virtualHost")
	flag.StringVar(&exchange, "exchange", "my.exchage", "RabbitMQ exchange")
	flag.StringVar(&exchangeRoutingKey, "exchangeRoutingKey", "my.routeKey", "RabbitMQ exchange routing key")
	flag.StringVar(&exchangeType, "exchangeType", "topic", "RabbitMQ exchange type")
	flag.BoolVar(&exchangeDurable, "exchangeDurable", true, "RabbitMQ exchange durability")
	flag.BoolVar(&exchangeAutoDelete, "exchangeAutoDelete", false, "RabbitMQ exchange auto-delete")

	flag.StringVar(&message, "message", "Hello Fucked World!", "Message to send into the queue")
	flag.StringVar(&messageContentType, "messageContentType", "application/json", "Message to send into the queue")

	flag.Parse()

	// Connecting
	uri := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, serverAddress, serverPort)

	amqpConfig := amqp.Config{}

	if virtualHost != "" {
		amqpConfig.Vhost = virtualHost
	}

	log.Printf("Connection to: %s", uri)
	conn, err := amqp.DialConfig(uri, amqpConfig)
	if err != nil {
		log.Fatalf("Failed to connec to to RabbitMQ: %s", err)
	}
	defer conn.Close()
	defer log.Println("Closing connection")

	// Go for Channel
	log.Println("Creating a channel")
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open RabbitMQ Channel: %s", err)
	}
	defer ch.Close()
	defer log.Println("Closing channel")

	err = ch.ExchangeDeclare(
		exchange,           // name
		exchangeType,       // type
		exchangeDurable,    // durable
		exchangeAutoDelete, // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Fatalf("Failed to Declare RabbitMQ Exchange: %s", err)
	}

	// Publishing
	log.Printf("Publishing message: %s, into the exchange: %s", message, exchange)
	msgConf := amqp.Publishing{}
	msgConf.ContentType = messageContentType
	msgConf.Body = []byte(message)
	if exchangeDurable {
		msgConf.DeliveryMode = amqp.Persistent
	}

	if err := ch.Publish(
		exchange,
		exchangeRoutingKey,
		false, // mandatory
		false, // immediate
		msgConf,
	); err != nil {
		log.Fatalf("Failed to publish a message on RabbitMQ exchange: %s", err)
	}

	log.Println("Message sent...")
}