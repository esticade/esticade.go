package esticade

func NewService(serciceName string) Service {
	service := new(esticadeService)
	service.name = serciceName
	return service
}

type esticadeService struct {
	name string
}

func (service *esticadeService) On(eventName string, callback func()) {

}

func (service *esticadeService) AlwaysOn(eventName string, callback func()) {

}

func (service *esticadeService) Emit(eventName string, payload interface{}) {

}

func (service *esticadeService) EmitChain(eventName string, payload interface{}) {

}

func (service *esticadeService) Shutdown() {

}