package util

import (
	log "github.com/sirupsen/logrus"
)

func HandleError(message string, err error) {
	if err != nil {
		log.Errorf("%s: %s", message, err)
		panic(err)
	}
}
