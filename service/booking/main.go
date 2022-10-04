package main

import (
	"Microservices/lib/config"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue/amqp"
	"Microservices/lib/util"
	"Microservices/service/booking/listener"
	"Microservices/service/booking/rest"
	"flag"
)

func main() {
	configPath := flag.String("conf", "./res/booking/configuration.json", "Configuration filepath")
	flag.Parse()

	conf := util.PackReturn(config.LoadConfiguration(*configPath))

	eventListener := util.PackReturnExit(amqp.NewAMQPListener("booking", "book", conf.AMQPMessageBroker))
	dbService := util.PackReturnExit(db.NewDatabaseService(db.DB_MONGODB, conf.DatabaseConnection))

	processor := listener.EventProcessor{
		EventListener:   eventListener,
		DatabaseService: dbService,
	}

	// Process events concurrently
	go func() {
		util.LogError(processor.ProcessEvents())
	}()

	rest.ServeAPI(dbService, nil, &conf)
}
