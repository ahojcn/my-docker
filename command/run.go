package command

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

func Run(command string, tty bool) {
	cmd := exec.Command(command)
	// for kinds of namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	/**
	 * 	Start() 不会阻塞，所以需要用 Wait()
	 *	Run() 会阻塞
	 */
	if err := cmd.Start(); err != nil {
		log.Fatalln("Run Start error", err)
	}
	cmd.Wait()
}
