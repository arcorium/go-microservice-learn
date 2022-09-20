package rest

import (
	"Microservices/event/config"
	"Microservices/event/persistence/db"
	"Microservices/event/util"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func CreateRouter(service db.DatabaseService) *mux.Router {
	handler := NewEventServiceHandler(service)

	router := mux.NewRouter().StrictSlash(false)
	eventRouter := router.PathPrefix("/events").Subrouter()
	eventRouter.Path("/{type}/{value}").Methods(http.MethodGet).HandlerFunc(handler.ServeHTTP)
	eventRouter.Path("").Methods(http.MethodGet).HandlerFunc(handler.FindAllEvent)
	eventRouter.Path("").Methods(http.MethodPost).HandlerFunc(handler.CreateEvent)

	return router
}

func ServeAPI(config_ *config.ServiceConfig, service db.DatabaseService) <-chan error {
	router := CreateRouter(service)
	router.Path("/").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		util.GetReturn(writer.Write([]byte("Index")))
	})

	log.Println("Server listening to", config_.Endpoint)

	// Channel for handling returned error
	httpErrorChan := make(chan error)

	go func(channel_ chan<- error) { channel_ <- http.ListenAndServe(config_.Endpoint, router) }(httpErrorChan)

	return httpErrorChan
}

func TLSServeAPI(config_ *config.HttpsConfig, service db.DatabaseService) <-chan error {
	router := CreateRouter(service)
	router.Path("/").Methods(http.MethodGet).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		util.GetReturn(writer.Write([]byte("Index TLS")))
	})

	log.Println("TLS Server listening to", config_.Endpoint)

	httpsErrorChan := make(chan error)

	go func(channel_ chan<- error, router_ *mux.Router) {
		channel_ <- http.ListenAndServeTLS(config_.Endpoint,
			config_.CertificatePath, config_.KeyPath, router_)
	}(httpsErrorChan, router)

	return httpsErrorChan
}
