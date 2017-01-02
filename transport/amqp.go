package transport


type AmqpService interface {
	Connect() error
	Emit(eventName, exchangeKey string, payload []byte) error
	On(eventName, exchangeKey string, callback func(body []byte) error) error
	Shutdown() error
}