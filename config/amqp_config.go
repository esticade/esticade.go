package config

import (
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const AMQP_DEFAULT_URL = "amqp://guest:guest@localhost:5672/"
const AMQP_DEFAULT_EXCHANGE_NAME = "events"
const AMQP_DEFAULT_ENGRAVED_VALUE = false

const ESTICADE_CONNECTION_URL_ENV_PROP = "ESTICADE_CONNECTION_URL"
const ESTICADE_EXCHANGE_NAME_ENV_PROP = "ESTICADE_EXCHANGE"
const ESTICADE_ENGRAVED_ENV_PROP = "ESTICADE_ENGRAVED"

const ESTICADE_CUSTOM_CONFIG_FILE_LOCATION = "ESTICADERC"
const ESTICADE_CONFIG_FILE_NAME = "esticade.json"

type amqpConfig struct {
	AmqpUrl      string `json:"connectionURL"`
	ExchangeName string `json:"exchange"`
	Engraved     bool `json:"engraved"`
}

func GetAmqpConfig() AmqpConfig {
	if confFile := os.Getenv(ESTICADE_CUSTOM_CONFIG_FILE_LOCATION); confFile != "" {
		if config, err := readConfigFromFile(confFile); err == nil {
			return config
		}
	}
	workDir, _ := os.Getwd()
	if config, err := readConfigFromFile(workDir + "/" + ESTICADE_CONFIG_FILE_NAME); err == nil {
		return config
	}
	return amqpConfig{AmqpUrl: AMQP_DEFAULT_URL, ExchangeName: AMQP_DEFAULT_EXCHANGE_NAME, Engraved: AMQP_DEFAULT_ENGRAVED_VALUE}
}

func (config amqpConfig) GetAmqpUrl() string {
	if url := os.Getenv(ESTICADE_CONNECTION_URL_ENV_PROP); url != "" {
		return url
	}
	return config.AmqpUrl;
}

func (config amqpConfig) GetExchangeName() string {
	if exchange := os.Getenv(ESTICADE_EXCHANGE_NAME_ENV_PROP); exchange != "" {
		return exchange
	}
	return config.ExchangeName;
}

func (config amqpConfig) GetEngraved() bool {
	if engraved := os.Getenv(ESTICADE_ENGRAVED_ENV_PROP); engraved != "" {
		bool, _ := strconv.ParseBool(engraved)
		return bool
	}
	return config.Engraved;
}

func readConfigFromFile(configFile string) (*amqpConfig, error) {
	config := new(amqpConfig)
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, config)
	return config, err
}
