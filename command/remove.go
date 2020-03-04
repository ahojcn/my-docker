package command

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
)

func Remove(containerName string)  {
	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorf("GetContainerInfo error:%v", err)
		return
	}
	if containerInfo.Status != container.STOP {
		log.Errorf("Could not remove not stopped container!")
		return
	}
	RemoveContainerInfo(containerInfo)
}
