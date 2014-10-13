package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type SyslogFields struct {
	timestamp, host, program, pid, message string
}

func extractFields(messageMap map[string]interface{}, fieldMapping FieldMapping) SyslogFields {
	formattedTimestamp := convertJSONTimestamp(messageMap[fieldMapping.Timestamp].(string))

	host := messageMap[fieldMapping.Host].(string)
	program := messageMap[fieldMapping.Program].(string)
	pid := fmt.Sprintf("[%s]", messageMap[fieldMapping.Pid].(string))
	message := messageMap[fieldMapping.Message].(string)

	return SyslogFields{
		formattedTimestamp,
		host,
		program,
		pid,
		message,
	}
}

func parseLogMessage(message []byte) map[string]interface{} {
	reader := bufio.NewReader(bytes.NewReader(message))
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	//skip first line containing the target index
	scanner.Scan()
	scanner.Scan()

	var msgif interface{}
	text := scanner.Text()
	err := json.Unmarshal([]byte(text), &msgif)
	failOnError(err, "Unmarshalling json failed")

	return msgif.(map[string]interface{})
}

func convertJSONTimestamp(timestr string) string {
	timestamp, err := time.Parse("2006-01-02T15:04:05.000-07:00", timestr)
	failOnError(err, "Timestamp parsing failed.")
	return timestamp.Format("Jan _2 15:04:05")
}
