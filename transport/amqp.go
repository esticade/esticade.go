package transport

import (
	"github.com/streadway/amqp"
	"fmt"
)

func NewAmqp() *transportAmqp {
	return NewAmqpByConfig(ConfigAmqp{})
}

func NewAmqpByConfig(config ConfigAmqp) *transportAmqp {
	transport := new(transportAmqp)

	transport.config = config
	if transport.config.Username == "" {
		transport.config.Username = "guest"
	}
	if transport.config.Password == "" {
		transport.config.Password = "guest"
	}
	if transport.config.Host == "" {
		transport.config.Host = "localhost"
	}
	if transport.config.Port == 0 {
		transport.config.Port = 5672
	}

	return transport
}

type ConfigAmqp struct {
	Username string
	Password string
	Host string
	Port int
}

type transportAmqp struct {
	config     ConfigAmqp
	connection *amqp.Connection
	channel *amqp.Channel
}

func (self *transportAmqp) Connect() error {
	var err error

	self.connection, err = amqp.Dial(self.formatUrl())
	if err != nil {
		return fmt.Errorf("Connection opening problem. %s", err.Error())
	}

	self.channel, err = self.connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel opening problem. %s", err.Error())
	}

	return nil
}

func (self *transportAmqp) Disconnect() error {
	var err error

	if self.channel != nil {
		err = self.channel.Close()
		if err != nil {
			return fmt.Errorf("Connection closing problem. %s", err.Error())
		}
	}

	if self.connection != nil {
		err = self.connection.Close()
		if err != nil {
			return fmt.Errorf("Channel closing problem. %s", err.Error())
		}
	}

	return nil
}

func (self *transportAmqp) formatUrl() string {
	return fmt.Sprintf(
		"amqp://%s%s:%d/",
		self.formatAuthString(),
		self.config.Host,
		self.config.Port,
	)
}

func (self *transportAmqp) formatAuthString() string {
	if self.config.Username== "" && self.config.Password == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s@", self.config.Username, self.config.Password)
}