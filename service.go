package esticade

import "github.com/esticade/esticade.go/transport"

func NewService(serviceName string) (Service, error) {
	return NewServiceByTransport(serviceName, transport.NewAmqp())
}

func NewServiceByTransport(serciceName string, transport Transport) (Service, error) {
	service := new(esticadeService)
	service.name = serciceName
	service.transport = transport

	errorConnect := service.transport.Connect()
	if errorConnect != nil {
		return nil, errorConnect
	}

	return service, nil
}

type esticadeService struct {
	name string
	transport Transport
}

func (self *esticadeService) On(eventName string, callback func()) {

}

func (self *esticadeService) AlwaysOn(eventName string, callback func()) {

}

func (self *esticadeService) Emit(eventName string, payload interface{}) {

}

func (self *esticadeService) EmitChain(eventName string, payload interface{}) {

}

func (self *esticadeService) Shutdown() error {
	return self.transport.Disconnect()
}