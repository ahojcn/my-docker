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

func Run(command string, tty bool, cg *cgroups.CgroupManger, rootPath string, volumes []string, containerName string) {
	reader, writer, err := os.Pipe()
	if err != nil {
		log.Errorln("os.Pope() error:", err)
		return
	}
	// cmd := exec.Command("/proc/self/exe", "init")

	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		fmt.Errorf("get init process error %v", err)
		return
	}
	cmd := exec.Command(initCmd, "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	newRootPath := getRootPath(rootPath)
	cmd.Dir = newRootPath + "/busybox"
	if err := NewWorkDir(newRootPath, volumes); err == nil {
		cmd.Dir = newRootPath + "/mnt"
	}
	defer ClearWorkDir(newRootPath, volumes)

	cmd.ExtraFiles = []*os.File{reader}
	sendInitCommand(command, writer)

	id := ContainerUUID()
	if containerName == "" {
		containerName = id
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		logFile, err := GetLogFile(containerName)
		if err != nil {
			log.Errorln("get log file error:", err)
			return
		}
		cmd.Stdout = logFile
	}

	if err := cmd.Start(); err != nil {
		log.Fatalln("Run cmd.Start error", err)
	}

	log.Infof("before process pid:%d, memory limit: %s", cmd.Process.Pid, cg.SubsystemsIns)

	cg.Set()
	defer cg.Destroy()
	cg.Apply(strconv.Itoa(cmd.Process.Pid))

	RecordContainerInfo(strconv.Itoa(cmd.Process.Pid), containerName, id, command)

	if tty {
		cmd.Wait()
		DeleteContainerInfo(containerName)
	}
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
		log.Infof("rootPath is empty, set cmd.Dir by default: %sbusybox", defaultPath)
		rootPath = defaultPath
	}
	imageTar := rootPath + "/busybox.tar"
	exist, _ := PathExists(imageTar)
	if !exist {
		log.Warnf("%s does not exist, set cmd.Dir by default: %sbusybox", imageTar, defaultPath)
		return defaultPath
	}
	imagePath := rootPath + "/busybox"
	exist, _ = PathExists(imageTar)
	if exist {
		os.RemoveAll(imagePath)
	}
	if err := os.Mkdir(imagePath, 0777); err != nil {
		log.Warnf("mkdir %s error: %v, set cmd.Dir by default: %sbusybox", imagePath, err, defaultPath)
		return defaultPath
	}
	if _, err := exec.Command("tar", "-xvf", imageTar, "-C", imagePath).CombinedOutput(); err != nil {
		log.Warnf("tar -xvf %s -c %s, err: %v, set cmd.Dir by default: %sbusybox", imageTar, imagePath, err, defaultPath)
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
func NewWorkDir(rootPath string, volumes []string) error {
	if err := CreateContainerLayer(rootPath); err != nil {
		return fmt.Errorf("create container layer %s error: %v", rootPath, err)
	}
	if err := CreateMntPoint(rootPath); err != nil {
		return fmt.Errorf("create mnt point %s error: %v", rootPath, err)
	}
	if err := SetMountPoint(rootPath); err != nil {
		return fmt.Errorf("set mount point %s error: %v", rootPath, err)
	}

	for _, volume := range volumes {
		if err := CreateVolume(rootPath, volume); err != nil {
			return fmt.Errorf("create volume %s error: %v", volume, err)
		}
	}

	return nil
}

// 生成 rootPath/writeLayer 文件夹
func CreateContainerLayer(rootPath string) error {
	writerLayer := rootPath + "/writeLayer"
	if err := os.Mkdir(writerLayer, 0777); err != nil {
		log.Warnf("mkdir %s error:%v", writerLayer, err)
		return fmt.Errorf("mkdir %s error:%v", writerLayer, err)
	}
	return nil
}

// 生成 rootPath/mnt 文件夹
func CreateMntPoint(rootPath string) error {
	mnt := rootPath + "/mnt"
	if err := os.Mkdir(mnt, 0777); err != nil {
		log.Warnf("mkdir %s error:%v", mnt, err)
		return fmt.Errorf("mkdir %s error:%v", mnt, err)
	}
	return nil
}

// 挂载（比如：mount -t aufs -o dirs=/home/ahojcn/writeLayer:/home/ahojcn/busybox none /home/ahojcn/mnt）
func SetMountPoint(rootPath string) error {
	dirs := "dirs=" + rootPath + "/writeLayer:" + rootPath + "/busybox"
	mnt := rootPath + "/mnt"
	if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mnt).CombinedOutput(); err != nil {
		log.Errorf("mount -t aufs -o %s none %s, err:%v", dirs, mnt, err)
		return fmt.Errorf("mount -t aufs -o %s none %s, err:%v", dirs, mnt, err)
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
func ClearWorkDir(rootPath string, volumes []string) {
	for _, volume := range volumes {
		ClearVolume(rootPath, volume)
	}

	ClearMountPoint(rootPath)
	ClearWriteLayer(rootPath)
}

func ClearMountPoint(rootPath string) {
	mnt := rootPath + "/mnt"
	if _, err := exec.Command("umount", "-f", mnt).CombinedOutput(); err != nil {
		log.Errorf("umount -f %s error: %v", mnt, err)
	}
	if err := os.RemoveAll(mnt); err != nil {
		log.Errorf("remove %s error: %v", mnt, err)
	}
}

func ClearWriteLayer(rootPath string) {
	writeLayer := rootPath + "/writeLayer"
	if err := os.RemoveAll(writeLayer); err != nil {
		log.Errorf("remove %s error: %v", writeLayer, err)
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
				log.Errorf("mkdir %s error: %v", hostPath, err)
				return fmt.Errorf("mkdir %s error: %v", hostPath, err)
			}
		}
		mountPath := strings.Split(volume, ":")[1]
		containerPath := containerMntPath + mountPath
		if err := os.Mkdir(containerPath, 0777); err != nil {
			log.Errorf("mkdir %s error: %v", containerPath, err)
			return fmt.Errorf("mkdir %s error: %v", containerPath, err)
		}
		dirs := "dirs=" + hostPath
		if _, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerPath).CombinedOutput(); err != nil {
			log.Errorf("mount -t aufs -o %s none %s error: %v", dirs, containerPath, err)
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
			log.Errorf("umount -f %s error: %v", containerPath, err)
		}
		if err := os.RemoveAll(containerPath); err != nil {
			log.Errorf("remove %s error: %v", containerPath, err)
		}
	}
}
