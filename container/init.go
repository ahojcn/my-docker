package container

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

/*
容器初始化操作
这里的 init 函数是在容器内部执行的
也就是说代码执行到这里后，容器所在的进程其实就已经创建出来了，这是本容器执行的第一个进程。
使用 mount 先去挂在 proc 文件系统，以便后面通过 ps 等系统命令去查看当前进程资源的情况。
*/
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	setUpMount()

	// 调用 LookPath 可以在系统的 PATH 里面寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}

	log.Infof("找到了 path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}
	return nil
}

//func RunContainerInitProcess(command string, args []string) error {
//	logrus.Infof("command %s", command)
//
//	/*
//		MS_NOEXEC 在本文件系统中不允许执行其他程序
//		MS_NOSUID 在本系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID
//		MS_NODEV 自从 Linux 2.4 以来，所有 mount 的系统都会默认设定的参数。
//	*/
//	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
//	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
//	argv := []string{command}
//	/* syscall.Exec() 黑魔法   execve 系统调用
//	最终调用了 kernel 的 int execve(const char *filename, char *const argv[], char *const envp[]);
//	*/
//	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
//		logrus.Errorf(err.Error())
//	}
//	return nil
//}

func readUserCommand() []string {
	// uintptr(3) 就是指 index 为 3 的文件描述符，也就是传递进来的管道的一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

/**
Init 挂载点
*/
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("获取当前路径错误 %v", err)
		return
	}
	log.Infof("当前路径是 %s", pwd)

	// syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	pivotRoot(pwd)
	// see https://github.com/xianlubird/mydocker/issues/41
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	/**
	  为了使当前root的老root和新root不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
