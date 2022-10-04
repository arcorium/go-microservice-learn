package amqp

import (
	"Microservices/lib/model"
	"Microservices/lib/queue"
	"Microservices/lib/util"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue      string
	consumer   string
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer func(channel_ *amqp.Channel) {
		util.LogError(channel_.Close())
	}(channel)

	// Declaring queue for this service
	_, err = channel.QueueDeclare(a.queue, true,
		false, false, false, nil)

	return err
}

func (a *amqpEventListener) Listen(event ...string) (<-chan queue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	// Binding this queue to any key in events exchange
	for _, eventName := range event {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}

	// Consume queue
	recv, err := channel.Consume(a.queue, a.consumer, false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	eventChannel := make(chan queue.Event)
	errorChannel := make(chan error)

	// Run asynchronous for receiving 'delivery' and convert into Event
	go func() {
		// Defer in goroutine because it will run forever to convert Delivery to Event type
		defer func(channel_ *amqp.Channel) {
			util.LogError(channel_.Close())
		}(channel)

		for msg := range recv {
			// Construct Event
			eventName, err := msg.Headers["x-event-name"]
			if !err {
				util.LogError(msg.Nack(false, true))
				errorChannel <- errors.New("there is no x-event-name key in headers")
				continue
			}

			var ev queue.Event
			// Check event name and create struct
			switch eventName {
			// TODO: Change from hardcoded into splendid one
			case "event.created":
				ev = new(model.EventCreateEvent)

				if err := json.Unmarshal(msg.Body, &ev); err != nil {
					// Failed to unmarshall
					util.LogError(msg.Nack(false, true))
					errorChannel <- errors.New("failed to unmarshal")
					continue
				}

			default:
				util.LogError(msg.Nack(false, false))
				errorChannel <- errors.New("route-key is not listed")
				continue
			}
			util.LogError(msg.Ack(false))
			// Pass into another channel
			eventChannel <- ev
		}
	}()

	return eventChannel, errorChannel, nil
}

func NewAMQPListener(consumer string, queue string, uri string) (queue.EventListener, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	listener := amqpEventListener{
		connection: conn,
		queue:      queue,
		consumer:   consumer,
	}

	err = listener.setup()
	return &listener, err
}
