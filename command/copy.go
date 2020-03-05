package command

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"strings"
)

func Copy(src, dst string) {
	f1 := strings.Contains(src, ":")
	f2 := strings.Contains(dst, ":")
	if (f1 && f2) || (!f1 && !f2) {
		log.Errorf("f1: %v, f2: %v, not correct format!", f1, f2)
		return
	}

	var from_container_to_host bool = true
	containerUrl := src
	hostUrl := dst
	if f2 {
		from_container_to_host = false
		containerUrl = dst
		hostUrl = src
	}
	containerName := strings.Split(containerUrl, ":")[0]
	containerPath := strings.Split(containerUrl, ":")[1]
	log.Infof("containerUrl:%s, hostUrl:%s, containerName:%s, containerPath:%s", containerUrl, hostUrl, containerName, containerPath)

	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorln("Get Container Info error: ", err)
		return
	}
	containerMntPath := containerInfo.RootPath + "/mnt/" + containerName + containerPath
	hostPath := hostUrl
	log.Infof("containerMntPath:%s, hostPath:%s", containerMntPath, hostPath)
	log.Infoln("from_container_to_host:", from_container_to_host)

	if from_container_to_host {
		log.Infof("from %s to %s", containerMntPath, hostPath)
		FileCopy(containerMntPath, hostPath)
	} else {
		log.Infof("from %s to %s", hostPath, containerMntPath)
		FileCopy(hostPath, containerMntPath)
	}
}

func FileCopy(src, dst string) {
	exist, _ := PathExists(src)
	if !exist {
		log.Printf("src:%s not exists!\n", src)
		return
	}
	exist, _ = PathExists(dst)
	if !exist {
		log.Printf("dst:%s not exists!\n", src)
		return
	}
	if _, err := exec.Command("cp", "-r", src, dst).CombinedOutput(); err != nil {
		log.Printf("cp -r %s %s, err:%v\n", src, dst, err)
		return
	}
}
