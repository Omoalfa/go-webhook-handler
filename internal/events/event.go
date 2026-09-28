package events

type CreateEvent struct {
	Provider  string
	EventType string
	Payload   []byte
}
