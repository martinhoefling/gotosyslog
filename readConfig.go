package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type RegexMap map[][]string
type IncludeList RegexMap
type ExcludeList RegexMap

type FilterConfig struct {
	Include IncludeList
	Exclude ExcludeList
}

func defaultFilterConfig() interface{} {
	var c FilterConfig
	c.Include = make(map[][]string)
	c.Exclude = make(map[][]string)
	return c
}

func readFilterConfig(m map[string]interface{}) interface{} {
	var c FilterConfig
	c.Exclude = readRegexMap(m, "exclude", []string{})
	c.Include = readRegexMap(m, "include", []string{})
	return c
}

type ExchangeConfig struct {
	Name        string
	Type        string
	Durable     bool
	AutoDeleted bool
	Internal    bool
	NoWait      bool
}

func defaultExchangeConfig() interface{} {
	var c ExchangeConfig
	log.Print("Using Defaults for RabbitMQ Exchange Config")
	c.Name = "logs_input"
	c.Type = "direct"
	c.Durable = true
	c.AutoDeleted = false
	c.Internal = false
	c.NoWait = false
	return c
}

func readExchangeConfig(m map[string]interface{}) interface{} {
	var c ExchangeConfig
	log.Print("Reading RabbitMQ Exchange Config")
	c.Name = readString(m, "name", "logs_input")
	c.Type = readString(m, "type", "direct")
	c.Durable = readBool(m, "durable", true)
	c.AutoDeleted = readBool(m, "auto-deleted", false)
	c.Internal = readBool(m, "internal", false)
	c.NoWait = readBool(m, "no-wait", false)
	return c
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
}

func defaultQueueConfig() interface{} {
	var c QueueConfig
	log.Print("Using Defaults for RabbitMQ Queue Config")
	c.Name = "logs_queue"
	c.Durable = true
	c.AutoDelete = false
	c.Exclusive = false
	c.NoWait = false
	return c
}

func readQueueConfig(m map[string]interface{}) interface{} {
	var c QueueConfig
	log.Print("Reading RabbitMQ Queue Config")
	c.Name = readString(m, "name", "logs_queue")
	c.Durable = readBool(m, "durable", true)
	c.AutoDelete = readBool(m, "auto-delete", false)
	c.Exclusive = readBool(m, "exclusive", false)
	c.NoWait = readBool(m, "no-wait", false)
	return c
}

type BindingConfig struct {
	RoutingKey string
}

func defaultBindingConfig() interface{} {
	var c BindingConfig
	log.Print("Using Defaults for RabbitMQ Binding Config")
	c.RoutingKey = "log-forward"
	return c
}

func readBindingConfig(m map[string]interface{}) interface{} {
	var c BindingConfig
	log.Print("Reading RabbitMQ Binding Config")
	c.RoutingKey = readString(m, "routing-key", "log-forward")
	return c
}

type RabbitMQConfig struct {
	URL      string
	Exchange ExchangeConfig
	Queue    QueueConfig
	Binding  BindingConfig
}

func defaultRabbitMQConfig() interface{} {
	var c RabbitMQConfig
	log.Print("Using Defaults for RabbitMQ Config")
	c.URL = "amqp://guest:guest@localhost:5672/"
	c.Exchange = defaultExchangeConfig().(ExchangeConfig)
	c.Queue = defaultQueueConfig().(QueueConfig)
	c.Binding = defaultBindingConfig().(BindingConfig)
	return c
}

func readRabbitMQConfig(m map[string]interface{}) interface{} {
	var c RabbitMQConfig
	log.Print("Reading RabbitMQ Config")
	c.URL = readString(m, "url", "amqp://guest:guest@localhost:5672/")
	c.Exchange = readMap(m, "exchange", readExchangeConfig, defaultExchangeConfig).(ExchangeConfig)
	c.Queue = readMap(m, "queue", readQueueConfig, defaultQueueConfig).(QueueConfig)
	c.Binding = readMap(m, "binding", readBindingConfig, defaultBindingConfig).(BindingConfig)
	return c
}

type FieldMapping struct {
	Timestamp string
	Message   string
	Host      string
	Pid       string
	Program   string
}

func defaultFieldMapping() interface{} {
	var c FieldMapping
	log.Print("Using Defaults for FieldMapping Config")
	c.Timestamp = "@timestamp"
	c.Message = "message"
	c.Host = "host"
	c.Pid = "pid"
	c.Program = "program"
	return c
}

func readFieldMapping(m map[string]interface{}) interface{} {
	var c FieldMapping
	log.Print("Reading FieldMapping Config")
	c.Timestamp = readString(m, "timestamp", "@timestamp")
	c.Message = readString(m, "message", "message")
	c.Host = readString(m, "host", "host")
	c.Pid = readString(m, "pid", "pid")
	c.Program = readString(m, "program", "program")
	return c
}

type Config struct {
	RabbitMQ     RabbitMQConfig
	FieldMapping FieldMapping
	Filters      FilterConfig
	Output       string
}

func (c *Config) UnmarshalJSON(b []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(b, &m)
	log.Print("Reading Config")
	c.Output = readString(m, "output", "/etc/gotosyslog/config.json")
	c.RabbitMQ = readMap(m, "rabbitmq", readRabbitMQConfig, defaultRabbitMQConfig).(RabbitMQConfig)
	c.FieldMapping = readMap(m, "fieldmap", readFieldMapping, defaultFieldMapping).(FieldMapping)
	c.Filters = readMap(m, "filters", readFilterConfig, defaultFilterConfig).(FilterConfig)
	return err
}

func readConfig(filename string) (config Config) {
	jsonConfig, err := ioutil.ReadFile(filename)
	failOnError(err, "Config file could not be read")

	err = json.Unmarshal(jsonConfig, &config)
	failOnError(err, "Config file content could not be parsed")

	log.Printf("Config file %s read", filename)
	return
}
