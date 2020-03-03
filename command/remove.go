package command

import (
	log "github.com/Sirupsen/logrus"
)

func Remove(containerName string)  {
	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorf("GetContainerInfo error:%v", err)
		return
	}
	if containerInfo.Status != STOP {
		log.Errorf("Could not remove not stopped container!")
		return
	}
	RemoveContainerInfo(containerInfo)
}
