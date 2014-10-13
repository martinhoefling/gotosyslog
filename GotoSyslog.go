package main

import (
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func messageToLogLines(config Config, msgs <-chan amqp.Delivery, loglines chan<- string) {
	for msg := range msgs {
		loglines <- formatSyslogMessage(msg, config.FieldMapping)
	}
}

func writeLogLines(fileHandle *os.File, loglines <-chan string) {
	for logline := range loglines {
		fileHandle.WriteString(logline)
	}
}

func main() {
	configFilePtr := flag.String("c", "/etc/gotosyslog/config.json", "configuration file (json)")
	flag.Parse()
	config := readConfig(*configFilePtr)

	msgs, closeQueueFunction := setupQueueConsumer(config.RabbitMQ)
	defer closeQueueFunction()

	signals := make(chan os.Signal, 1)
	loglines := make(chan string, 100)

	fh, closeFhFunction := openOutfile(config)
	defer closeFhFunction()

	go messageToLogLines(config, msgs, loglines)
	go writeLogLines(fh, loglines)

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")

	signal.Notify(signals,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	<-signals
}
