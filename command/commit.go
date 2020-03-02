package command

import (
	log "github.com/Sirupsen/logrus"
	"os/exec"
)

func Commit(imageName string) {
	mntPath := DEFAULTPATH + "/mnt"
	imageTar := DEFAULTPATH + "/" + imageName + ".tar"
	log.Infoln("imageTar:", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		log.Errorf("tar -czf %s, -C %s error: %v\n", imageTar, mntPath, err)
	}
}
