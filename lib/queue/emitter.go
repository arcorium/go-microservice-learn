package queue

type Event interface {
	EventName() string
}

type EventEmitter interface {
	Emit(Event) error
}
