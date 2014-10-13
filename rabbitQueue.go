package main

import (
	"github.com/streadway/amqp"
	"log"
)

func setupQueueConsumer(config RabbitMQConfig) (<-chan amqp.Delivery, func()) {
	log.Printf("Connecting to RabbitMQ %s", config.URL)
	conn, err := amqp.Dial(config.URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	closeFunction := func() {
		log.Printf("Shutting down RabbitMQ connection")
		ch.Close()
		conn.Close()
	}

	err = ch.ExchangeDeclare(
		config.Exchange.Name,
		config.Exchange.Type,
		config.Exchange.Durable,
		config.Exchange.AutoDeleted,
		false, // internal
		config.Exchange.NoWait,
		nil, // arguments
	)
	failOnError(err, "Failed to declare exchange")

	q, err := ch.QueueDeclare(
		config.Queue.Name,
		config.Queue.Durable,
		config.Queue.AutoDelete,
		config.Queue.Exclusive,
		config.Queue.NoWait,
		nil, // arguments
	)
	failOnError(err, "Failed to declare queue")

	s := "log-forward"
	log.Printf("Binding queue %s to exchange %s with routing key %s", q.Name, "logs_input", s)
	err = ch.QueueBind(
		config.Queue.Name,
		config.Binding.RoutingKey,
		config.Exchange.Name,
		false,
		nil)
	failOnError(err, "Failed to bind queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register consumer")
	return msgs, closeFunction
}
