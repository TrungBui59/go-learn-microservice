package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/event"
	"log"
	"math"
	"time"
)

func main() {
	// try to connect to rabbitmq
	connection, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to RabbitMQ")
	defer connection.Close()
	// start to listening for messages
	log.Println("Listening and consuming RabbitMQ messages...")
	// create consumer
	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitMQ:5672/")
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %s", err)
			counts++
		} else {
			connection = c
			break
		}

		if counts >= 5 {
			log.Printf("Cannot connect after 5 retries")
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Sleeping %v", backOff)
		time.Sleep(backOff)
	}
	return connection, nil
}
