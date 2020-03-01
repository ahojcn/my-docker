package main

import (
	"mydocker/cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
	"os"
	"strconv"
	"testing"
	"time"
)

func Test003(t *testing.T) {
	subsystems.Set("10M")
	pid := os.Getpid()
	log.Infoln("current pid:", pid)
	subsystems.Apply(strconv.Itoa(pid))
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 1)
	}
}
