package esticade

import (
	"github.com/satori/go.uuid"
	"github.com/esticade/esticade.go/config"
	"github.com/esticade/esticade.go/transport"
	"encoding/json"
	"errors"
)

type amqpService struct {
	serviceName      string
	correlationBlock string
	config           config.AmqpConfig
	transportService transport.AmqpService
}

func NewService(serviceName string) (Service, error) {
	config := config.GetAmqpConfig()
	return NewServiceWithConfig(serviceName, config.GetAmqpUrl(), config.GetExchangeName(), config.GetEngraved())
}

func NewServiceWithConfig(serviceName, amqpUrl, exchangeName string, engraved bool) (Service, error) {
	service := &amqpService{
		serviceName: serviceName,
		correlationBlock: string(uuid.NewV4().Bytes()),
		transportService: transport.NewRabbitMqService(amqpUrl, exchangeName, engraved),
	}
	return service, service.transportService.Connect()
}

func (service *amqpService) Emit(eventName string, payload interface{}) error {
	payloadByte, err := json.Marshal(Event{EventId: string(uuid.NewV4().Bytes()), Name: eventName, Body: payload})
	if err != nil {
		return errors.New("Failed to encode payload:" + err.Error())
	}
	return service.transportService.Emit(service.correlationBlock + "." + eventName, payloadByte)
}

func (service *amqpService) On(eventName string, callback func(event Event)) error {
	return service.transportService.On(
		service.serviceName + "-" + eventName,
		"*." + eventName,
		func(body []byte) error {
			var event Event
			if err := json.Unmarshal(body, &event); err != nil {
				return err
			}
			callback(event)
			return nil
		},
	)
}

func (service *amqpService) Shutdown() error {
	return service.transportService.Shutdown()
}