package command

import (
	log "github.com/Sirupsen/logrus"
	"mydocker/cgroups"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func Run(command string, tty bool, cg *cgroups.CgroupManger, rootPath string) {
	// cmd := exec.Command(command)
	/**
	 * 先执行当前进程的 init 命令，参数为 command
	 *
	 * 表示在执行用户的 command 命令以前，在已经做好 namespace 隔离的进程中先执行 init 命令
	 * 其实就是执行 mount -t proc proc /proc 操作，然后再执行 command
	 */
	//cmd := exec.Command("/proc/self/exe", "init", command)
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Errorln("os.Pope() error:", err)
		return
	}
	cmd := exec.Command("/proc/self/exe", "init")

	// for kinds of namespace
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}

	log.Infoln("cmd.Dir:", "/root")
	// cmd.Dir = "/root"
	cmd.Dir = rootPath
	if rootPath == "" {
		log.Infoln("set cmd.Dir by default: /root/busybox")
		cmd.Dir = "/root/busybox"
	}
	cmd.ExtraFiles = []*os.File{reader}
	sendInitCommand(command, writer)

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

	// 直接设置 memory 限制
	log.Infof("--- before process pid:%d, memory limit: %s ---", cmd.Process.Pid, cg.SubsystemsIns)
	//subsystems.Set(memory)
	//subsystems.Apply(strconv.Itoa(cmd.Process.Pid))
	//defer subsystems.Remove()
	cg.Set()
	defer cg.Destroy()
	cg.Apply(strconv.Itoa(cmd.Process.Pid))

	cmd.Wait()
}

func sendInitCommand(command string, writer *os.File) {
	_, err := writer.Write([]byte(command))
	if err != nil {
		log.Errorln("write.Write error:", err)
		return
	}
	writer.Close()
}
