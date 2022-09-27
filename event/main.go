package main

import (
	"Microservices/event/config"
	"Microservices/event/persistence/db"
	"Microservices/event/rest"
	"Microservices/event/util"
	"Microservices/lib/queue/amqp"
	"flag"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	// load configuration file from user input or default value
	configPath := flag.String("conf", "./res/event/configuration.json", "Path to configuration")
	flag.Parse()

	// Load Configuration
	conf, _ := config.LoadConfiguration(*configPath)
	// Create broker or event emitter service
	broker := util.PackReturnExit(amqp.NewAMQPEventEmitter(conf.AMQPMessageBroker))
	// Create database service
	if svc, err := db.NewDatabaseService(conf.DatabaseType, conf.DatabaseConnection); err != nil {
		log.Println(err)
		return
	} else {
		// Run api both for http and https
		res := rest.ServeAPI(&conf, svc, broker)

		for _, channel := range res {
			// Wait until all receive-channel notified
			select {
			case err := <-channel:
				log.Println(err)
			}
		}
	}

	//amqpUri := os.Getenv("RABBITMQ_URI")
	//if amqpUri == "" {
	//	amqpUri = "amqp://guest:guest@localhost:5672"
	//}
	//if conn, err := amqp.Dial(amqpUri); err == nil {
	//	defer func(conn *amqp.Connection) {
	//		if err := conn.Close(); err != nil {
	//			log.Println(err)
	//		}
	//	}(conn)
	//
	//	// Create channel
	//	channel, err := conn.Channel()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	// durable: save the exchange even the rabbitmq is restarted
	//	// autoDelete: delete exchange when the channel is closed
	//	// internal: prevent publisher to publish message into this exchange
	//	// noWait: don't wait the response from broker and just return
	//	util.LogError(channel.ExchangeDeclare("events", "topic", true,
	//		false, false, false, nil))
	//
	//	// Declare queue
	//	util.PackReturnExit(channel.QueueDeclare("my_queue",
	//		true, false, false, false, nil))
	//	// Bind queue to exchange
	//	util.LogError(channel.QueueBind("my_queue", "topic",
	//		"events", false, nil))
	//
	//	// Publishing message
	//	message := amqp.Publishing{Body: []byte("Hellaw")}
	//	// mandatory: the message will publish into at least one queue
	//	// immediate: the message will publish into at least one subscriber
	//	// key will be used to route
	//	util.LogError(channel.PublishWithContext(context.Background(),
	//		"events", "topic", false, false, message))
	//
	//} else {
	//	log.Fatalln(err)
	//}
}
