package config

type AmqpConfig interface {
	GetAmqpUrl() string
	GetExchangeName() string
	GetEngraved() bool
}
