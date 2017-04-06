package esticade

type Service interface {
	Connect() error
	On(eventName string, callback func(event Event)) error
	Emit(eventName string, payload interface{}) error
	Shutdown() error
}

type Event struct {
	EventId string `json:"id"`
	Name    string `json:"name"`
	Body    interface{} `json:"body"`
}

func (event Event) GetEventId() string {
	return event.EventId
}

func (event Event) GetName() string {
	return event.Name
}

func (event Event) GetBody() interface{} {
	return event.Body
}