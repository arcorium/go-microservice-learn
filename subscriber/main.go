package main

import (
	"Microservices/event/util"
	_ "github.com/joho/godotenv/autoload"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"os"
)

func main() {
	amqpUri := os.Getenv("RABBITMQ_URI")
	if amqpUri == "" {
		amqpUri = "amqp://guest:guest@localhost:5672"
	}

	// Connect to amqp
	conn := util.PackReturnExit(amqp.Dial(amqpUri))
	// Create channel for asynchronous
	channel := util.PackReturn(conn.Channel())

	// autoAck: acknowledge received message automatically
	// exclusive: only this consumer should consume queue
	// noLocal: don't deliver message on the same channel
	message := util.PackReturn(channel.Consume("my_queue",
		"subscriber", false, false,
		false, false, nil))

	// Waiting channel to received message continuously using range instead of for and select
	for data := range message {
		log.Println(string(data.Body))
		// Acknowledge received message, because autoAck set to false
		// If Ack is not done, the message will still in queue
		util.LogError(data.Ack(false))
	}
}
