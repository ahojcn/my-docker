package main

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strings"
)

/*
这里的 Start 方法是真正开始前面创建好的 command 的调用
它首先会 clone 出来一个 namespace 隔离的进程
然后在子进程中，调用 /proc/self/exe，也就是调用自己，
	发送 init 参数，调用我们写的 init 方法初始化容器的一些资源
*/
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	log.Infoln("创建 Namespace 隔离的容器进程。")
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	// use mydocker-cgroup as cgroup name
	// 创建 cgroup manager，并通过调用 set 和 apply 设置资源限制并使限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	// 设置资源限制
	cgroupManager.Set(res)
	// 将容器进程加入到各个 subsystem 挂在对应的 cgroup 中
	cgroupManager.Apply(parent.Process.Pid)
	// 对容器设置完限制之后，初始化容器（发送用户命令）
	sendInitCommand(comArray, writePipe)
	parent.Wait()
	mntURL := "/root/mnt/"
	rootURL := "/root/"
	container.DeleteWorkSpace(rootURL, mntURL)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
