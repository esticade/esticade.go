package esticade

import (
	"encoding/json"
	"errors"
	"github.com/esticade/esticade.go/config"
	"github.com/esticade/esticade.go/transport"
	"github.com/satori/go.uuid"
)

type amqpService struct {
	serviceName      string
	correlationBlock string
	config           config.AmqpConfig
	transportService transport.AmqpService
}

func NewService(serviceName string) Service {
	config := config.GetAmqpConfig()
	return NewServiceCustomConfiguration(serviceName, config.GetAmqpUrl(), config.GetExchangeName(), config.GetEngraved())
}

func NewServiceCustomConfiguration(serviceName, amqpUrl, exchangeName string, engraved bool) Service {
	service := &amqpService{
		serviceName:      serviceName,
		correlationBlock: string(uuid.NewV4().Bytes()),
		transportService: transport.NewRabbitMqService(amqpUrl, exchangeName, engraved),
	}
	return service
}

func (service *amqpService) Connect() error {
	return service.transportService.Connect()
}

func (service *amqpService) Emit(eventName string, payload interface{}) error {
	payloadByte, err := json.Marshal(Event{EventId: string(uuid.NewV4().Bytes()), Name: eventName, Body: payload})
	if err != nil {
		return errors.New("Failed to encode payload:" + err.Error())
	}
	return service.transportService.Emit(service.correlationBlock+"."+eventName, payloadByte)
}

func (service *amqpService) On(eventName string, callback func(event Event) error) error {
	return service.transportService.On(
		service.serviceName+"-"+eventName,
		"*."+eventName,
		func(body []byte) error {
			var event Event
			if err := json.Unmarshal(body, &event); err != nil {
				return err
			}
			return callback(event)
		},
	)
}

func (service *amqpService) Shutdown() error {
	return service.transportService.Shutdown()
}
