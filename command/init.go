package command

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"syscall"
)

func Init(command string) {
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

	/**
	exec 函数族
	syscall.Exec 会先执行参数指定的命令，但是并不创建新的进程，只在当前进程空间内执行
		即替换当前进程的执行内容，会重用同一个进程PID
	*/
	if err := syscall.Exec(command, []string{command}, os.Environ()); err != nil {
		log.Fatalln("Init syscall.Exec() error", err)
	}
}
