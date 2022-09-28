package main

import (
	"Microservices/lib/config"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue/amqp"
	"Microservices/lib/util"
	"Microservices/services/booking/listener"
	"flag"
)

func main() {
	configPath := flag.String("conf", "./res/event/configuration.json", "Configuration filepath")
	flag.Parse()

	conf := util.PackReturn(config.LoadConfiguration(*configPath))

	eventListener := util.PackReturnExit(amqp.NewAMQPListener("booking", "book", conf.AMQPMessageBroker))
	dbService := util.PackReturnExit(db.NewDatabaseService(db.DB_MONGODB, conf.DatabaseConnection))

	processor := util.PackReturnExit(listener.NewEventProcessor(eventListener, dbService))

	util.LogError(processor.ProcessEvents())
}
