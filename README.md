# esticade.go
Simple event driven micro-services framework http://esticade.io/

# Installing 
Install using go get
``` bash
go get github.com/esticade/esticade.go
```

# API

To start using Esticade, please create an instance of esticade Service
object. 

## Service object

- `esticade.NewService(serviceName)` - Will construct a new service and connect to the exchange.                            
- `esticade.NewServiceWithConfig(serviceName, amqpUrl, exchangeName string, engraved bool)` - Will construct a new service with 
the given configuration and connect to the exchange.                            
- `on(eventName, callback)` - Will register event listener. Callback will be called with an `Event` object as the only 
argument. If there are two or more instances of the same service running, the events will be equally divided between all the instances. 
- `emit(eventName[, payload])` - Will emit event to the event network. Returns promise that is fulfilled once the event is emitted.
- `shutdown()` - Will shut the entire service down, if there is nothing else keeping process alive, the process will terminate.

## Event object

- `name` - Name of the event.
- `body` - The content of the message.
- `correlationId` - Will be same on all the events in the event chain.
- `eventId` - Unique identifier for the event


# Quick start

Install and run RabbitMQ and try out the hello esticade demo

Hello, esticade demo:
```go
package main

import (
	"github.com/esticade/esticade.go"
	"time"
	"log"
)

func main() {
	service, err := esticade.NewService("Esticade test");
	if err != nil {
		panic(err.Error())
	}
	if err = service.On("my-event", callback); err != nil {
		panic(err.Error())
	}
	if err = service.Emit("my-event", "Hello esticade"); err != nil {
		panic(err.Error())
	}

	time.Sleep(300 * time.Millisecond)
	service.Shutdown()
}

func callback(event esticade.Event) {
	log.Printf("received: %s, %s", event.GetName(), event.GetBody())
}
```

You should see the output:
```
received: my-event, Hello esticade
```

# Configuration

Using `esticade.NewService(serviceName)` by default connects to localhost with user and pass guest/guest. This is the default configuration
for RabbitMQ. If you want to override that, you can override it with a configuration file in any of following locations or use the 
`esticade.NewServiceWithConfig(serviceName, amqpUrl, exchangeName string, engraved bool)` constructor

- Environment variables for each of the configuration variables
- A file pointed to by ESTICADERC environment variable
- esticade.json in current working folder.

If any of those files is located in that order, it's contents are read and used for configuration. It should contain
JSON object with any of the following properties: 

- `connectionURL` - Default: `amqp://guest:guest@localhost/`
- `exchange` - Default `events`
- `engraved` - Default `false`. Will make named queues (those registered with service.on()) durable. We suggest you leave this
option to `false` during development as otherwise you will get a lot of permanent queues in the rabbitmq server. You should
turn this on in production though, as it will make sure no messages get lost when service restarts. Turning it off when it
has been turned on might cause errors as the durable queues are not redefined as non-durable automatically. You have
to manually delete the queues from RabbitMQ.

Example:

```json
{ 
    "connectionURL": "amqp://user:pass@example.com/vhost",
    "exchange": "EventNetwork"
}
```

## Environment variables

- `ESTICADE_CONNECTION_URL` - AMQP url to connect to
- `ESTICADE_EXCHANGE` - Exchange name
- `ESTICADE_ENGRAVED` - Whether or not to engrave the queues 
