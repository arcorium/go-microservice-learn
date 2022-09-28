package amqp

import (
	"Microservices/lib/queue"
	"Microservices/lib/util"
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
	defer func(channel_ *amqp.Channel) {
		util.LogError(channel_.Close())
	}(channel)

	// Declare exchange
	// durable: save the exchange even the rabbitmq is restarted
	// autoDelete: delete exchange when the channel is closed
	// internal: prevent publisher to publish message into this exchange
	// noWait: don't wait the response from broker and just return
	return channel.ExchangeDeclare("events", "topic", true,
		false, false, false, nil)
}

func (a amqpEventEmitter) Emit(event_ queue.Event) error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer func(channel_ *amqp.Channel) {
		util.LogError(channel_.Close())
	}(channel)
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Marshalling
	data, err := json.Marshal(event_)
	if err != nil {
		return err
	}

	// Construct message
	msg := amqp.Publishing{
		Headers:     amqp.Table{"x-event-name": event_.EventName()},
		Body:        data,
		ContentType: "application/json",
	}

	// Publish
	// mandatory: the message will publish into at least one queue
	// immediate: the message will publish into at least one services
	// key will be used to route
	return channel.PublishWithContext(ctx, "events",
		event_.EventName(), false, false, msg)
}

func NewAMQPEventEmitter(uri string) (queue.EventEmitter, error) {
	// Connect into AMQP
	connection, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	
	emitter := amqpEventEmitter{
		connection: connection,
	}
	// Setup emitter
	err = emitter.setup()

	return emitter, err
}
