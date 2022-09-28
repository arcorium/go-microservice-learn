package rest

import (
	"Microservices/lib/config"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue"
	"Microservices/lib/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateRouter(service db.DatabaseService, emitter_ queue.EventEmitter) *mux.Router {
	handler := NewEventServiceHandler(service, emitter_)

	router := mux.NewRouter().StrictSlash(false)
	eventRouter := router.PathPrefix("/events").Subrouter()
	eventRouter.Path("/{type}/{value}").Methods(http.MethodGet).HandlerFunc(handler.ServeHTTP)
	eventRouter.Path("").Methods(http.MethodGet).HandlerFunc(handler.FindAllEvent)
	eventRouter.Path("").Methods(http.MethodPost).HandlerFunc(handler.CreateEvent)

	return router
}

func serve(config_ *config.ServiceConfig, service db.DatabaseService, emitter_ queue.EventEmitter) <-chan error {
	router := CreateRouter(service, emitter_)
	router.Path("/").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		util.PackReturn(writer.Write([]byte("Index")))
	})

	log.Println("Server listening to", config_.Endpoint)
	log.Println("Broker connected to", config_.AMQPMessageBroker)

	// Channel for handling returned error
	httpErrorChan := make(chan error)

	go func(channel_ chan<- error) { channel_ <- http.ListenAndServe(config_.Endpoint, router) }(httpErrorChan)

	return httpErrorChan
}

func serveTLS(config_ *config.HttpsConfig, service db.DatabaseService, emitter_ queue.EventEmitter) <-chan error {
	router := CreateRouter(service, emitter_)
	router.Path("/").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		util.PackReturn(writer.Write([]byte("Index TLS")))
	})

	log.Println("TLS Server listening to", config_.Endpoint)

	httpsErrorChan := make(chan error)

	go func(channel_ chan<- error, router_ *mux.Router) {
		channel_ <- http.ListenAndServeTLS(config_.Endpoint,
			config_.CertificatePath, config_.KeyPath, router_)
	}(httpsErrorChan, router)

	return httpsErrorChan
}

func ServeAPI(config_ *config.ServiceConfig, service db.DatabaseService, emitter_ queue.EventEmitter) []<-chan error {
	httpChan := serve(config_, service, emitter_)

	result := make([]<-chan error, 0)
	result = append(result, httpChan)
	// Serve for https
	if config_.Https {
		httpsChan := serveTLS(&config_.HttpsConfig, service, emitter_)
		result = append(result, httpsChan)
	}

	return result
}
