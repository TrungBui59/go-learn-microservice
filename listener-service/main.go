package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
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

	// create consumer

	// watch the queue and consume events
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// dont continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
