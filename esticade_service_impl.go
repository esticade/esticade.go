package esticade

import (
	"github.com/satori/go.uuid"
	"github.com/esticade.go/config"
	"github.com/esticade.go/transport"
	"encoding/json"
	"errors"
)

type amqpService struct {
	serviceName      string
	correlationBlock string
	config           config.AmqpConfig
	transportService transport.AmqpService
}

type callbackAdapter struct {
	callback func(event Event)
}

func (adapter callbackAdapter) doCallback(body []byte) error {
	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}
	adapter.callback(event)
	return nil
}

func NewService(serviceName string) (Service, error) {
	config := config.GetAmqpConfig()
	service := &amqpService{
		serviceName: serviceName,
		correlationBlock: string(uuid.NewV4().Bytes()),
		transportService: transport.NewRabbitMqService(config.GetAmqpUrl(), config.GetExchangeName(), config.GetEngraved()),
	}
	return service, service.transportService.Connect()
}

func (service *amqpService) Emit(eventName string, payload interface{}) error {
	payloadByte, err := json.Marshal(Event{EventId: string(uuid.NewV4().Bytes()), Name: eventName, Body: payload})
	if err != nil {
		return errors.New("Failed to encode payload:" + err.Error())
	}
	return service.transportService.Emit(eventName, service.correlationBlock + "." + eventName, payloadByte)
}

func (service *amqpService) On(eventName string, callback func(event Event)) error {
	return service.transportService.On(eventName, "*." + eventName, callbackAdapter{callback}.doCallback)
}

func (service *amqpService) Shutdown() error {
	return service.transportService.Shutdown()
}