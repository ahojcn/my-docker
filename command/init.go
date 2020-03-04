package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
)

func Init(command string) {
	log.Infoln("read from commandline:", command)
	command = readFromPipe()
	log.Infoln("read from pipe:", command)

	pwd, err := os.Getwd()
	if err != nil {
		log.Errorln("get pwd error:", err)
		return
	}
	log.Infoln("current path:", pwd)
	pivotRoot(pwd)

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

func readFromPipe() string {
	reader := os.NewFile(uintptr(3), "pipe")
	command, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Errorln("reader.Read(buf) error:", err)
		return ""
	}
	return string(command)
}

/**
把路径 root 变为根节点
*/
func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount root fs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的 rootfs, 现在老的 old_root 是挂载在 rootfs/.pivot_root
	// 挂载点现在依然可以在 mount 命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root: %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir error: %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
