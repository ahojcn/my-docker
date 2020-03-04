package main

import (
	"mydocker/cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
	"testing"
)

func Test000(t *testing.T) {
	mountPath := subsystems.FindCgroupMountPoint("memory")
	log.Infoln("mountPath:", mountPath)
}

func Test001(t *testing.T) {
	absolutePath := subsystems.FindAbsolutePath("memory")
	log.Infoln("absolutePath:", absolutePath)
}
