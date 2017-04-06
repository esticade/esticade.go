package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/esticade/esticade.go/config"
	"testing"
	"os"
	"time"
	"strings"
	"github.com/esticade/esticade.go"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration suite")
}

var _ = Describe("Rabbit MQ service", func() {

	var (
		serviceInstance esticade.Service
		eventName string
	)

	BeforeEach(func() {
		os.Setenv(ESTICADE_CONNECTION_URL_ENV_PROP, AMQP_DEFAULT_URL)
		os.Setenv(ESTICADE_EXCHANGE_NAME_ENV_PROP, "events_integration_test")
		os.Setenv(ESTICADE_ENGRAVED_ENV_PROP, "false")
		serviceInstance = esticade.NewService("Test service")
		serviceInstance.Connect()
		eventName = "test_event"
	})

	AfterEach(func() {
		serviceInstance.Shutdown()
	})

	Context("Published and subscribes to messages", func() {

		It("Publishes events with Emit", func() {
			callbackCounter := consumerInvocationCounter{invokedCount: 0}
			serviceInstance.On(eventName, callbackCounter.call)

			serviceInstance.Emit(eventName, "testMsg")
			serviceInstance.Emit(eventName, 1)
			serviceInstance.Emit(eventName, true)

			time.Sleep(500 * time.Millisecond)
			Expect(callbackCounter.invokedCount).To(Equal(3))
		})

		It("Invokes correct subscribers on event", func() {
			errorSubscriber := consumerInvocationCounter{invokedCount: 0}
			secretInfoSubscriber := consumerInvocationCounter{invokedCount: 0}

			serviceInstance.On("ERROR", errorSubscriber.call)
			serviceInstance.On("SECRET_INFO", secretInfoSubscriber.call)

			serviceInstance.Emit("ERROR", "testMsg")
			serviceInstance.Emit("ERROR", 1)
			serviceInstance.Emit("ERROR", true)

			serviceInstance.Emit("SECRET_INFO", "8675309 - Jenny")

			time.Sleep(500 * time.Millisecond)

			Expect(errorSubscriber.invokedCount).To(Equal(3))
			Expect(secretInfoSubscriber.invokedCount).To(Equal(1))
		})
	})
	It("Shuts down the connection when Shutdown invoked", func() {

		errorSubscriber := consumerInvocationCounter{invokedCount: 0}
		serviceInstance.On("something", errorSubscriber.call)
		serviceInstance.Shutdown()
		time.Sleep(300 * time.Millisecond)
		error := serviceInstance.Emit(eventName, 1)
		Expect(strings.HasSuffix(error.Error(), "Exception (504) Reason: \"channel/connection is not open\"")).To(Equal(true))

	})

	Context("Handles different message payload types", func() {

		It("Supports strings", func() {
			verifyInvoked("Hello world!", serviceInstance, eventName)
		})

		It("Supports numbers", func() {
			verifyInvoked(200.23, serviceInstance, eventName)
		})

		It("Supports booleans", func() {
			verifyInvoked(true, serviceInstance, eventName)
		})
	})

})

func verifyInvoked(message interface{}, serviceInstance esticade.Service, eventName string) {
	callbackCounter := consumerInvocationCounter{invokedCount: 0, expectedMessage: message}
	serviceInstance.On(eventName, callbackCounter.callIfBodyEqualsExpected)
	serviceInstance.Emit(eventName, message)
	time.Sleep(1 * time.Second)
	Expect(callbackCounter.invokedCount).To(Equal(1))
}

type consumerInvocationCounter struct {
	invokedCount    int
	expectedMessage interface{}
}

func (counter *consumerInvocationCounter) callIfBodyEqualsExpected(event esticade.Event) {
	if counter.expectedMessage == event.Body {
		counter.invokedCount = counter.invokedCount + 1
	}
}

func (counter *consumerInvocationCounter) call(event esticade.Event) {
	counter.invokedCount = counter.invokedCount + 1

}