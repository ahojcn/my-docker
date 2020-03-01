package command

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"syscall"
)

func Init(command string) {
	log.Infoln("read from commandline:", command)
	command = readFromPipe()
	log.Infoln("read from pipe:", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	/**
	// 5-2 直接 cmd.Run() 会导致 1 号进程不是用户进程
	cmd := exec.Command(command)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalln("Init Run() error", err)
	}
	*/

	pwd, err := os.Getwd()
	if err != nil {
		log.Errorln("get pwd error:", err)
		return
	}
	log.Infoln("current path:", pwd)

	/**
	exec 函数族
	syscall.Exec 会先执行参数指定的命令，但是并不创建新的进程，只在当前进程空间内执行
		即替换当前进程的执行内容，会重用同一个进程PID
	*/
	if err := syscall.Exec(command, []string{command}, os.Environ()); err != nil {
		log.Fatalln("Init syscall.Exec() error", err)
	}
}

func readFromPipe() string {
	reader := os.NewFile(uintptr(3), "pipe")
	command, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorln("reader.Read(buf) error:", err)
		return ""
	}
	return string(command)
}
