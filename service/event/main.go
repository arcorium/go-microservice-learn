package main

import (
	"Microservices/lib/config"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue/amqp"
	"Microservices/lib/util"
	"Microservices/service/event/rest"
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
}
