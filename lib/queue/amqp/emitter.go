package amqp

import (
	"Microservices/event/util"
	"Microservices/lib/queue"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type amqpEventEmitter struct {
	connection *amqp.Connection
}

func (a amqpEventEmitter) setup() error {
	// Create channel
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer func(channel *amqp.Channel) {
		util.LogError(channel.Close())
	}(channel)

	// Declare exchange
	return channel.ExchangeDeclare("events", "topic", true,
		false, false, false, nil)
}

func (a amqpEventEmitter) Emit(event_ queue.Event) error {
	// Create new channel
	// channel created every emitting event to be usable in multithreading
	channel := util.PackReturn(a.connection.Channel())
	defer func(channel *amqp.Channel) {
		util.LogError(channel.Close())
	}(channel)
	// Marshalling
	data := util.PackReturn(json.Marshal(event_))
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Construct message
	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event_.EventName()},
		Body:        data,
		ContentType: "application/json",
	}
	// Publish
	return channel.PublishWithContext(ctx, "events",
		event_.EventName(), false, false, msg)
}

func NewAMQPEventEmitter(uri string) (queue.EventEmitter, error) {
	// Connect into AMQP
	connection, err := amqp.Dial(uri)

	emitter := amqpEventEmitter{
		connection: connection,
	}
	// Setup emitter
	err = emitter.setup()

	return emitter, err
}
