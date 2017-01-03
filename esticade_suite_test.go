package esticade

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/satori/go.uuid"
	"github.com/golang/mock/gomock"
	"github.com/esticade.go/testing/mocks"
)

func TestEsticade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("Service interface", func() {
	eventName := "event_name"
	eventBody := "event_text"
	var (
		mockCtrl *gomock.Controller
		transportFake *mocks.MockAmqpService
		uuidString string
		service *amqpService
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		transportFake = mocks.NewMockAmqpService(mockCtrl)
		uuidString = string(uuid.NewV4().Bytes())
		service = &amqpService{
			serviceName: "Name",
			correlationBlock: uuidString,
			transportService: transportFake,
		}

	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Context("When emiting new event", func() {
		It("Calls Emit on amqp transport with correct parameters", func() {
			transportFake.EXPECT().Emit(eventName, uuidString + "." + eventName, gomock.Any())
			service.Emit(eventName, eventBody)
		})
	})
	Context("Adding new listener", func() {
		It("Registers listener with amqp transport", func() {
			transportFake.EXPECT().On(eventName, "*." + eventName, gomock.Any())
			service.On(eventName, testFunc)
		})
	})
	Context("Initiating shutdown", func() {
		It("Calls shutdown on amqp transport", func() {
			transportFake.EXPECT().Shutdown()
			service.Shutdown()
		})
	})
})

func testFunc(event Event) {

}