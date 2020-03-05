package main

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"testing"
)

func Test011(t *testing.T) {
	err := os.MkdirAll("/root/1111/2/3/4/5", os.ModePerm)
	log.Infoln(err)
}
