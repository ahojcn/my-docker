package command

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/container"
	"strconv"
	"syscall"
)

func Stop(containerName string) {
	log.Infoln("停止容器中")
	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorf("GetContainerInfo error:%v", err)
		return
	}
	if containerInfo.Pid == "" {
		log.Infoln("container not exists!")
		return
	}
	pid, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		log.Errorf("strconv.Atoi(%s) error : %v", containerInfo.Pid, err)
		return
	}
	if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
		log.Errorf("Stop container %s error : %v", containerName, err)
		return
	}
	containerInfo.Status = container.STOP
	containerInfo.Pid = ""
	log.Infoln(containerName, "已停止")
	UpdateContainerInfo(containerInfo)

	log.Infoln("rootPath:", containerInfo.RootPath)
	log.Infoln(containerInfo.Volumes)
	ClearWorkDir(containerInfo.RootPath, containerName, containerInfo.Volumes)
}
