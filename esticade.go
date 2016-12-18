package esticade

type Service interface {
	On(eventName string, callback func())
	AlwaysOn(eventName string, callback func())
	Emit(eventName string, payload interface{})
	EmitChain(eventName string, payload interface{})
	Shutdown()
}
