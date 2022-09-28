package queue

type EventListener interface {
	Listen(event ...string) (<-chan Event, <-chan error, error)
}
