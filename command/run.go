package command

import (
	"../cgroups/subsystems"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func Run(command string, tty bool, memory string) {
	// cmd := exec.Command(command)
	/**
	 * 先执行当前进程的 init 命令，参数为 command
	 *
	 * 表示在执行用户的 command 命令以前，在已经做好 namespace 隔离的进程中先执行 init 命令
	 * 其实就是执行 mount -t proc proc /proc 操作，然后再执行 command
	 */
	cmd := exec.Command("/proc/self/exe", "init", command)

	// for kinds of namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
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

	log.Infof("--- before process pid:%d, memory limit: %s ---", cmd.Process.Pid, memory)
	subsystems.Set(memory)
	subsystems.Apply(strconv.Itoa(cmd.Process.Pid))

	cmd.Wait()
}
