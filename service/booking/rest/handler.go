package rest

import (
	"Microservices/lib/config"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue"
	"Microservices/lib/util"
	"github.com/gorilla/mux"
	"net/http"
)

type BookingEventHandler struct {
	Db      db.DatabaseService
	Emitter queue.EventEmitter
}

func (b *BookingEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func ServeAPI(dbService db.DatabaseService, emitter queue.EventEmitter, configuration *config.ServiceConfig) {
	router := mux.NewRouter()
	router.Handle("/event/{eventId}/booking", &BookingEventHandler{Db: dbService, Emitter: emitter}).Methods(http.MethodPost)

	util.LogError(http.ListenAndServe(configuration.Endpoint, router))
}
