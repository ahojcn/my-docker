package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"mydocker/cgroups"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func Run(command string, tty bool, cg *cgroups.CgroupManger, rootPath string, volume string) {
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

	// log.Infoln("cmd.Dir:", "/root")
	// cmd.Dir = "/root"
	newRootPath := getRootPath(rootPath)
	cmd.Dir = newRootPath + "/busybox"
	if err := NewWorkDir(newRootPath, volume); err == nil {
		cmd.Dir = newRootPath + "/mnt"
	}
	defer ClearWorkDir(newRootPath, volume)
	//if rootPath == "" {
	//	log.Infoln("set cmd.Dir by default: /root/busybox")
	//	cmd.Dir = "/root/busybox"
	//}
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

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
根据命令行 -m 提供的目录返回执行 init 程序的目录。
比如输入 `./mydocker run -it -r /home/ahojcn /bin/sh`
	此时 rootPath = /home/ahojcn
	并且需要 /home/ahojcn 目录已经准备好了 busybox.tar 文件
	此时程序会解压 busybox.tar 到 /home/ahojcn/busybox 并把 /home/ahojcn/busybox 设置为执行 init 程序的工作目录
*/
const (
	DEFAULTPATH = "/home/ahojcn/test-mydocker/"
)

func getRootPath(rootPath string) string {
	log.Infoln("rootPath:", rootPath)
	defaultPath := DEFAULTPATH
	if rootPath == "" {
		log.Infof("rootPath is empty, set cmd.Dir by default: %s/busybox\n", defaultPath)
		rootPath = defaultPath
	}
	imageTar := rootPath + "/busybox.tar"
	exist, _ := PathExists(imageTar)
	if !exist {
		log.Warnf("%s does not exist, set cmd.Dir by default: %s/busybox\n", imageTar, defaultPath)
		return defaultPath
	}
	imagePath := rootPath + "/busybox"
	exist, _ = PathExists(imageTar)
	if exist {
		os.RemoveAll(imagePath)
	}
	if err := os.Mkdir(imagePath, 0777); err != nil {
		log.Warnf("mkdir %s error: %v, set cmd.Dir by default: %s/busybox\n", imagePath, err, defaultPath)
		return defaultPath
	}
	if _, err := exec.Command("tar", "-xvf", imageTar, "-C", imagePath).CombinedOutput(); err != nil {
		log.Warnf("tar -xvf %s -c %s, err: %v, set cmd.Dir by default: %s/busybox\n", imageTar, imagePath, err, defaultPath)
		return defaultPath
	}

	return rootPath
}

/**
创建挂载点工作
1. 创建 writeLayer 文件夹
2. 创建 mnt 文件夹
3. 挂载：将 busybox 和 writeLayer 挂载到 mnt 下。
*/
// 创建 Init 程序工作目录
func NewWorkDir(rootPath, volume string) error {
	if err := CreateContainerLayer(rootPath); err != nil {
		return fmt.Errorf("create container layer %s error: %v\n", rootPath, err)
	}
	if err := CreateMntPoint(rootPath); err != nil {
		return fmt.Errorf("create mnt point %s error: %v\n", rootPath, err)
	}
	if err := SetMountPoint(rootPath); err != nil {
		return fmt.Errorf("set mount point %s error: %v\n", rootPath, err)
	}

	if err := CreateVolume(rootPath, volume); err != nil {
		return fmt.Errorf("create volume %s error: %v\n", volume, err)
	}

	return nil
}

// 生成 rootPath/writeLayer 文件夹
func CreateContainerLayer(rootPath string) error {
	writerLayer := rootPath + "/writeLayer"
	if err := os.Mkdir(writerLayer, 0777); err != nil {
		log.Warnf("mkdir %s error:%v\n", writerLayer, err)
		return fmt.Errorf("mkdir %s error:%v\n", writerLayer, err)
	}
	return nil
}

// 生成 rootPath/mnt 文件夹
func CreateMntPoint(rootPath string) error {
	mnt := rootPath + "/mnt"
	if err := os.Mkdir(mnt, 0777); err != nil {
		log.Warnf("mkdir %s error:%v\n", mnt, err)
		return fmt.Errorf("mkdir %s error:%v\n", mnt, err)
	}
	return nil
}

// 挂载（比如：mount -t aufs -o dirs=/home/ahojcn/writeLayer:/home/ahojcn/busybox none /home/ahojcn/mnt）
func SetMountPoint(rootPath string) error {
	dirs := "dirs=" + rootPath + "/writeLayer:" + rootPath + "/busybox"
	mnt := rootPath + "/mnt"
	if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mnt).CombinedOutput(); err != nil {
		log.Errorf("mount -t aufs -o %s none %s, err:%v\n", dirs, mnt, err)
		return fmt.Errorf("mount -t aufs -o %s none %s, err:%v\n", dirs, mnt, err)
	}
	log.Warnln("mount success!")

	return nil
}

/**
清除挂载点工作
1. umount /home/ahojcn/mnt
2. rmdir /home/ahojcn/mnt
3. rmdir /home/ahojcn/writeLayer
*/
func ClearWorkDir(rootPath, volume string) {
	ClearVolume(rootPath, volume)


	ClearMountPoint(rootPath)
	ClearWriteLayer(rootPath)
}

func ClearMountPoint(rootPath string) {
	mnt := rootPath + "/mnt"
	if _, err := exec.Command("umount", "-f", mnt).CombinedOutput(); err != nil {
		log.Errorf("umount -f %s error: %v\n", mnt, err)
	}
	if err := os.RemoveAll(mnt); err != nil {
		log.Errorf("remove %s error: %v\n", mnt, err)
	}
}

func ClearWriteLayer(rootPath string) {
	writeLayer := rootPath + "/writeLayer"
	if err := os.RemoveAll(writeLayer); err != nil {
		log.Errorf("remove %s error: %v\n", writeLayer, err)
	}
}

/**
处理 volume 的添加和删除方法
*/
// 增加 volume 并且 mount
func CreateVolume(rootPath, volume string) error {
	if volume != "" {
		containerMntPath := rootPath + "/mnt"
		hostPath := strings.Split(volume, ":")[0]
		exist, _ := PathExists(hostPath)
		if !exist {
			if err := os.Mkdir(hostPath, 0777); err != nil {
				log.Errorf("mkdir %s error: %v\n", hostPath, err)
				return fmt.Errorf("mkdir %s error: %v", hostPath, err)
			}
		}
		mountPath := strings.Split(volume, ":")[1]
		containerPath := containerMntPath + mountPath
		if err := os.Mkdir(containerPath, 0777); err != nil {
			log.Errorf("mkdir %s error: %v", containerPath, err)
			return fmt.Errorf("mkdir %s error: %v\n", containerPath, err)
		}
		dirs := "dirs=" + hostPath
		if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerPath).CombinedOutput(); err != nil {
			log.Errorf("mount -t aufs -o %s none %s error: %v\n", dirs, containerPath, err)
			return fmt.Errorf("mount -t aufs -o %s none %s error: %v", dirs, containerPath, err)
		}
	}

	return nil
}

// 删除 volume
func ClearVolume(rootPath, volume string) {
	if volume != "" {
		containerMntPath := rootPath + "/mnt"
		mountPath := strings.Split(volume, ":")[1]
		containerPath := containerMntPath + mountPath
		if _, err := exec.Command("umount", "-f", containerPath).CombinedOutput(); err != nil {
			log.Errorf("umount -f %s error: %v\n", containerPath, err)
		}
		if err := os.RemoveAll(containerPath); err != nil {
			log.Errorf("remove %s error: %v\n", containerPath, err)
		}
	}
}
