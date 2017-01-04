package transport


type AmqpService interface {
	Connect() error
	Emit(exchangeKey string, payload []byte) error
	On(queueName, exchangeKey string, callback func(body []byte) error) error
	Shutdown() error
}