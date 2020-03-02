package subsystems

import (
	"bufio"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"strings"
)

/**找到对应 subsystem 的目录位置
根据 subsystem 的类型找到对应的 hierarchy
从而可以在该 hierarchy 创建子 cgroup
进而把进程添加到此 cgroup 的限制中
从而达到在此 subsystem 上限制进程的作用
*/
func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		log.Fatalln("Open mountinfo error:", err)
		return ""
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	for {
		line, err := bfRd.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return ""
			}
		}

		parts := strings.Split(string(line), " ")
		if strings.Contains(parts[len(parts)-1], subsystem) {
			//log.Warnln("parts[4]:", parts[4])
			return parts[4]
		}
	}
}

/**找到当前容器所在 subsystem 的 hierarchy 的绝对路径
根据 subsystem 需要找到当前容器的 cgroup 位置，这样才可以往里面加入相关的限制
*/
func FindAbsolutePath(subsystem string) string {
	path := FindCgroupMountPoint(subsystem)
	if path != "" {
		absolutePath := path + "/" + ResourceName
		exist, err := PathExists(absolutePath)
		if err != nil {
			log.Fatalln("Path exists error", err)
			return ""
		}
		if !exist {
			err := os.Mkdir(absolutePath, os.ModePerm)
			if err != nil {
				log.Fatalln("Mkdir absolutePath error:", err)
				return ""
			}
		}

		return absolutePath
	}
	return ""
}

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
