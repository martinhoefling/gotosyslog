package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func formatSyslogMessage(message amqp.Delivery, fieldMapping FieldMapping) string {
	msgmap := parseLogMessage(message.Body)
	fields := extractFields(msgmap, fieldMapping)
	return fmt.Sprintf("%s %s %s%s: %s\n", fields.timestamp, fields.host, fields.program, fields.pid, fields.message)
}
