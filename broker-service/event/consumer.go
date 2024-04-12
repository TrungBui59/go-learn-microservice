package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
)

// receiving end of the queue
type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	// setup the exchange for log
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil

}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	// get the channel
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	// close when done
	defer ch.Close()

	// get a queue
	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	// bind our chanel to each topic
	for _, topic := range topics {
		err = ch.QueueBind(q.Name, topic, "logs_topic", false, nil)
		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			err := json.Unmarshal(d.Body, &payload)
			if err != nil {
				log.Println(err)
			}
			log.Printf("Received a message: %s", payload)
			go handlePayload(payload)

		}
	}()

	fmt.Printf(" [*] Waiting for messages [Exchange, Queue] [log_topics, %s]. To exit press CTRL+C\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		//log whenever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		//authenticate
	default:
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logEvent(payload Payload) error {
	// get the data
	jsonData, _ := json.MarshalIndent(payload, "", "\t")

	logServiceURL := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
