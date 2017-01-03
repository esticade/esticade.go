package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/esticade.go/config"
	"testing"
	"os"
	"io/ioutil"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Reading amqp configuration", func() {

	Context("When environment properties are set", func() {
		os.Setenv(ESTICADE_CONNECTION_URL_ENV_PROP, "amqp://test:123")
		os.Setenv(ESTICADE_EXCHANGE_NAME_ENV_PROP, "events")
		os.Setenv(ESTICADE_ENGRAVED_ENV_PROP, "true")

		It("Should use the env. property values", func() {
			Expect(GetAmqpConfig().GetAmqpUrl()).To(Equal("amqp://test:123"))
			Expect(GetAmqpConfig().GetExchangeName()).To(Equal("events"))
			Expect(GetAmqpConfig().GetEngraved()).To(Equal(true))
			os.Unsetenv(ESTICADE_CONNECTION_URL_ENV_PROP)
			os.Unsetenv(ESTICADE_EXCHANGE_NAME_ENV_PROP)
			os.Unsetenv(ESTICADE_ENGRAVED_ENV_PROP)
		})
	})

	Context("When custom config file environment property is set", func() {
		os.Setenv(ESTICADE_CUSTOM_CONFIG_FILE_LOCATION, "../testing/resources/esticade.json")

		It("Should read configuration from custom config file", func() {
			Expect(GetAmqpConfig().GetAmqpUrl()).To(Equal("amqp://user:pass@example.com/vhost"))
			Expect(GetAmqpConfig().GetExchangeName()).To(Equal("test_exchange"))
			Expect(GetAmqpConfig().GetEngraved()).To(Equal(true))
			os.Unsetenv(ESTICADE_CUSTOM_CONFIG_FILE_LOCATION)
		})
	})

	Context("When esticade.json file provided in working directory", func() {
		jsonString := "{\"connectionURL\": \"amqp://user:pass@example.com/\", \"exchange\": \"test_exchange\", \"engraved\": true } "
		ioutil.WriteFile(ESTICADE_CONFIG_FILE_NAME, []byte(jsonString), 0644)
		It("Should read and use the esticade.json file", func() {
			Expect(GetAmqpConfig().GetAmqpUrl()).To(Equal("amqp://user:pass@example.com/"))
			Expect(GetAmqpConfig().GetExchangeName()).To(Equal("test_exchange"))
			Expect(GetAmqpConfig().GetEngraved()).To(Equal(true))
			os.Remove(ESTICADE_CONFIG_FILE_NAME)
		})
	})

	Context("When no configuration is provided", func() {
		It("Should use default configuration", func() {
			Expect(GetAmqpConfig().GetAmqpUrl()).To(Equal(AMQP_DEFAULT_URL))
			Expect(GetAmqpConfig().GetExchangeName()).To(Equal(AMQP_DEFAULT_EXCHANGE_NAME))
			Expect(GetAmqpConfig().GetEngraved()).To(Equal(AMQP_DEFAULT_ENGRAVED_VALUE))
		})
	})
})