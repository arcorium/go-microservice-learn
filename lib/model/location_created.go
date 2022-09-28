package model

import (
	"Microservices/services/event/model"
)

type LocationCreateEvent struct {
	Id      string       `json:"locationId"`
	Name    string       `json:"name"`
	Address string       `json:"address"`
	Country string       `json:"country"`
	Halls   []model.Hall `json:"halls,omitempty"`
}

func (l *LocationCreateEvent) EventName() string {
	return "location.created"
}
