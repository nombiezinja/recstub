package main

import (
	"encoding/json"
	"fmt"
	"log"

	e "github.com/nombiezinja/chstub/entities"

	u "github.com/nombiezinja/chstub/utils"
	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	u.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	u.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"outgoingStuffs", // name
		false,            // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	u.FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	u.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var p e.Payload
			err := json.Unmarshal(d.Body, &p)

			u.FailOnError(err, "Failed to unmarshall")
			fmt.Printf("Received a payload %+v: \n ", p)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
