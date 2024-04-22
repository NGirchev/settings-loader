package util

import (
	log "github.com/sirupsen/logrus"
	"sync"
)

func DoWork(id int, workFunc func(), wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Worker %d started", id)
	workFunc()
	log.Debugf("Worker %d stopped", id)
}
