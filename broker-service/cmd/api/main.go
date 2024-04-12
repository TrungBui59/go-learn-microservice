package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	// try to connect to rabbitmq
	connection, err := connect()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to RabbitMQ")
	defer connection.Close()

	app := Config{
		Rabbit: connection,
	}
	log.Printf("Starting broker service on port %s\n", webPort)

	//defin http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.route(),
	}

	// start to server
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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
