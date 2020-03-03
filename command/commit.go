package command

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

func Commit(containerName, imageName string) {
	// mntPath := DEFAULTPATH + "/mnt"
	// imageTar := DEFAULTPATH + "/" + imageName + ".tar"
	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorln("GetContainerInfo error:", err)
		return
	}
	mntPath := containerInfo.RootPath + "/mnt/" + containerName
	imageTar := containerInfo.RootPath + "/" + imageName + ".tar"
	log.Infoln("imageTar:", imageTar)
	log.Infoln("mntPath:", mntPath)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar -czf %s, -C %s error: %v\n", imageTar, mntPath, err)
	}
}
