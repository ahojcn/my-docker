package command

import log "github.com/Sirupsen/logrus"

func Logs(containerName string) {
	data := ReadLogs(containerName)
	log.Infoln(data)
}
