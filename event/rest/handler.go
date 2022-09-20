package rest

import (
	"Microservices/event/model"
	"Microservices/event/persistence/db"
	"Microservices/event/util"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type EventServiceHandler struct {
	dbService db.DatabaseService
}

func NewEventServiceHandler(dbService_ db.DatabaseService) EventServiceHandler {
	return EventServiceHandler{dbService: dbService_}
}

func (e EventServiceHandler) ServeHTTP(writer_ http.ResponseWriter, req_ *http.Request) {
	util.SetDefaultHeader(writer_)

	parameters := mux.Vars(req_)
	types := parameters["type"]
	values := parameters["value"]

	switch types {
	case "id":
		if res := e.findEventById(values); res == nil {
			writer_.WriteHeader(http.StatusBadRequest)
		} else {
			util.IsError(json.NewEncoder(writer_).Encode(*res))
		}
		break
	case "name":
		if res := e.findEventByName(values); res == nil {
			writer_.WriteHeader(http.StatusBadRequest)
		} else {
			util.IsError(json.NewEncoder(writer_).Encode(*res))
		}
		break
	default:
		writer_.WriteHeader(http.StatusBadRequest)
	}
}

func (e EventServiceHandler) findEventById(id_ any) *model.Event {
	if event, err := e.dbService.FindEventById(id_); err != nil {
		return nil
	} else {
		//fmt.Printf("Address of event in handler find by id : %p\n", event)
		return event
	}
}

func (e EventServiceHandler) findEventByName(name_ string) *model.Event {
	return util.GetReturn(e.dbService.FindEventByName(name_))
}

func (e EventServiceHandler) FindAllEvent(writer_ http.ResponseWriter, req_ *http.Request) {
	if events, err := e.dbService.FindAllEvents(); err != nil {
		writer_.WriteHeader(http.StatusBadRequest)
	} else {
		util.SetDefaultHeader(writer_)
		util.LogError(json.NewEncoder(writer_).Encode(&events))
	}
}

func (e EventServiceHandler) CreateEvent(writer_ http.ResponseWriter, req_ *http.Request) {
	var event model.Event

	if util.IsError(json.NewDecoder(req_.Body).Decode(&event)) {
		return
	}
	if id, err2 := e.dbService.AddEvent(&event); err2 != nil {
		log.Println(err2)
	} else {
		util.GetReturn(writer_.Write([]byte(id.(string))))
	}
}