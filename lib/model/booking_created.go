package model

type BookCreatedEvent struct {
	EventId string `json:"eventId"`
	UserId  string `json:"userId"`
}

func (b *BookCreatedEvent) EventName() string {
	return "event.booked"
}
