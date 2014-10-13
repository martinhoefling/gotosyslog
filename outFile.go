package main

import (
	"log"
	"os"
)

func openOutfile(config Config) (*os.File, func()) {
	fh, err := os.OpenFile(config.Output, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		fh, err = os.Create(config.Output)
	}

	failOnError(err, "Opening file for syslog output failed.")

	closeFunc := func() {
		log.Print("Syncing filehandle.")
		fh.Sync()
		log.Print("Closing filehandle.")
		fh.Close()
		return
	}
	return fh, closeFunc
}
