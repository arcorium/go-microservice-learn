package contract

import "time"

type EventCreateEvent struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	LocationId string    `json:"locationId"`
	TimeStart  time.Time `json:"timeStart"`
	TimeEnd    time.Time `json:"timeEnd"`
}

func (c *EventCreateEvent) EventName() string {
	return "event.created"
}
