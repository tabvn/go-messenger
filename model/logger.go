package model

import (
	"os"
	"log"
	"messenger/config"
)

func writeToLog(message string, kindOfMessage string) {

	f, err := os.OpenFile(config.WorkingDir+"/messenger.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, kindOfMessage, log.LstdFlags)
	logger.Println(message)
}
