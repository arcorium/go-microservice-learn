package main

import (
	"Microservices/event/config"
	"Microservices/event/persistence/db"
	"Microservices/event/rest"
	"flag"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {

	// Checking configuration
	configPath := flag.String("conf", "./res/event/configuration.json", "Path to configuration")
	flag.Parse()

	conf, _ := config.LoadConfiguration(*configPath)

	if svc, err := db.NewDatabaseService(conf.DatabaseType, conf.DatabaseConnection); err != nil {
		log.Println(err)
		return
	} else {
		// Run api both for http and https
		httpErrorChan := rest.ServeAPI(&conf, svc)

		// Run https server when it is configured
		if conf.Https {
			httpsErrorChan := rest.TLSServeAPI(&conf.HttpsConfig, svc)

			// Wait until https channel notified with error
			select {
			case https := <-httpsErrorChan:
				log.Println(https)
			}
		}

		// Wait until http channel notified with error
		select {
		case http := <-httpErrorChan:
			log.Println(http)
		}
	}

}
