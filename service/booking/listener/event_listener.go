package listener

import (
	"Microservices/lib/model"
	"Microservices/lib/persistence/db"
	"Microservices/lib/queue"
	"log"
)

type EventProcessor struct {
	EventListener   queue.EventListener
	DatabaseService db.DatabaseService
}

func (e EventProcessor) ProcessEvents() error {
	eventChan, errorChan, err := e.EventListener.Listen("book.created")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case event := <-eventChan:
			switch event.EventName() {
			case "book.created":
				book := event.(*model.BookCreatedEvent)
				log.Printf("Event %s booked by %s", book.EventId, book.UserId)
			}
		case err := <-errorChan:
			log.Println(err)
		}
	}
}
