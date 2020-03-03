package command

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
)

func Exec(containerName, command string) {
	containerInfo, err := GetContainerInfo(containerName)
	if err != nil {
		log.Errorln("get container info error:", err)
		return
	}
	pid := containerInfo.Pid
	os.Setenv("mydocker_pid", pid)
	os.Setenv("mydocker_cmd", command)

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error %v\n", containerName, err)
	}
}
