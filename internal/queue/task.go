package queue

const (
	TaskProcessEvent = "event:process"
)

type EventTaskPayload struct {
	EventID string `json:"event_id"`
}
