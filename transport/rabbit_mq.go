package transport

import (
	"github.com/streadway/amqp"
	"fmt"
	"errors"
	"time"
	"log"
)

const MESSAGE_CONTENT_TYPE_JSON = "application/json"

type rabbitMqService struct {
	correlationBlock    string
	url                 string
	exchange            string
	shuttingDown        bool
	durable             bool
	connection          *amqp.Connection
	connectionCloseChan chan *amqp.Error
	consumers           []*consumer
}

type consumer struct {
	eventName   string
	exchangeKey string
	callback    func(body []byte) error
}

func NewRabbitMqService(url, exchange string, durable bool) AmqpService {
	return &rabbitMqService{
		url: url,
		exchange: exchange,
		durable: durable,
		shuttingDown: false,
		connectionCloseChan: make(chan *amqp.Error),
	}
}
func (service *rabbitMqService) Connect() error {
	if service.connection == nil {
		if err := service.connect(); err != nil {
			return err
		}
	}
	return nil
}

func (service *rabbitMqService) Emit(eventName, exchangeKey string, payload []byte) error {

	channel, err := service.getChannel()
	if err != nil {
		return errors.New("Failed to open a channel:" + err.Error())
	}

	if err = createExchange(channel, service.exchange); err != nil {
		return fmt.Errorf("Failed to declare exchange: %s", err)
	}

	return channel.Publish(service.exchange, exchangeKey, false, false,
		amqp.Publishing{
			ContentType: MESSAGE_CONTENT_TYPE_JSON,
			Body: payload,
		})
}

func (service *rabbitMqService) On(eventName, exchangeKey string, callback func(body []byte) error) error {
	channel, err := service.getChannel()
	if err != nil {
		return errors.New("Failed to create a channel: " + err.Error())
	}

	queue, err := createQueue(channel, eventName, service.durable)
	if err != nil {
		return errors.New("Failed to declare a queue: " + err.Error())
	}

	if err = createExchange(channel, service.exchange); err != nil {
		return fmt.Errorf("Failed to declare exchange: %s", err)
	}

	if err = channel.QueueBind(queue.Name, exchangeKey, service.exchange, false, nil); err != nil {
		return errors.New("Failed to bind queue to exchange: " + err.Error())
	}

	deliveries, err := consumeQueue(channel, queue.Name)
	if err != nil {
		return errors.New("Failed to consume queue messages: " + err.Error())
	}

	go startEventHandler(deliveries, callback)
	service.consumers = append(service.consumers, &consumer{eventName: eventName, exchangeKey: exchangeKey, callback: callback})
	return nil
}

func (service *rabbitMqService) Shutdown() error {
	service.shuttingDown = true
	return service.connection.Close()
}

func (service *rabbitMqService) getChannel() (*amqp.Channel, error) {
	if service.connection == nil {
		if err := service.connect(); err != nil {
			return nil, err
		}
	}
	return service.connection.Channel()
}

func (service *rabbitMqService) connect() error {
	var err error
	service.connection, err = amqp.Dial(service.url)
	if err != nil {
		return errors.New("Failed to connect to server: " + err.Error())
	}
	go service.reconnectOnErrorHandler()
	return nil

}

func (service *rabbitMqService) reconnectOnErrorHandler() {
	var amqpError *amqp.Error
	service.connection.NotifyClose(service.connectionCloseChan)
	for {
		amqpError = <-service.connectionCloseChan
		log.Println("Received connection close event for " + service.url)
		if service.shuttingDown {
			return
		} else if amqpError != nil {
			service.connection = reconnectToServer(service.url)
			service.connectionCloseChan = make(chan *amqp.Error)
			service.connection.NotifyClose(service.connectionCloseChan)
			service.reattachConsumers()
		}
	}
}

func (service *rabbitMqService) reattachConsumers() {
	for _, element := range service.consumers {
		service.On(element.eventName, element.exchangeKey, element.callback)
	}
}

func startEventHandler(deliveries <-chan amqp.Delivery, callback func(body []byte) error) {
	for delivery := range deliveries {
		if err := callback(delivery.Body); err == nil {
			delivery.Ack(false)
		} else {
			log.Println("Error handling received event: " + err.Error())
		}
	}
}

func reconnectToServer(url string) *amqp.Connection {
	for {
		log.Println("Reconnecting to " + url)
		conn, err := amqp.Dial(url)
		if err == nil {
			return conn
		}
		log.Println(fmt.Sprintf("Failed to reconnect to %s : %s", url, err.Error()))
		time.Sleep(500 * time.Millisecond)
	}
}

func consumeQueue(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return channel.Consume(queueName, "", false, false, false, false, nil, )
}

func createQueue(channel *amqp.Channel, queueName string, durable bool) (amqp.Queue, error) {
	return channel.QueueDeclare(queueName, durable, false, false, false, nil)
}

func createExchange(channel *amqp.Channel, name string) error {
	return channel.ExchangeDeclare(name, amqp.ExchangeTopic, true, false, false, false, nil)
}