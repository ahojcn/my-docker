package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	containerEnvs := getEnvsByPid(containerInfo.Pid)
	cmd.Env = append(os.Environ(), containerEnvs...)

	if err := cmd.Run(); err != nil {
		log.Errorf("exec container %s error %v", containerName, err)
	}
}

/*
由于进程存放环境变量的位置是 /proc/<pid>/environ
因此根据给定 pid 读取这个文件，便可以获取环境变量。
在文件内容中，每个环境变量之间通过 \u0000 分割
*/
func getEnvsByPid(pid string) []string {
	path := fmt.Sprintf("/proc/%s/environ", pid)
	contentBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Infof("Read file %s error %v", path, err)
		return nil
	}
	envs := strings.Split(string(contentBytes), "\u0000")
	return envs
}
