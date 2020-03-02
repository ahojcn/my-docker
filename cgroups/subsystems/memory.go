package subsystems

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type MemorySubsystem struct {
}

func (s *MemorySubsystem) Name() string {
	return "memory"
}

/**
向 memory.limit_in_bytes 写入限制 content
*/
func (s *MemorySubsystem) Set(res *ResourceConfig) error {
	if res.MemoryLimit != "" {
		content := res.MemoryLimit
		absolutePath := ""
		if absolutePath = FindAbsolutePath(s.Name()); absolutePath == "" {
			log.Errorln("absolutePath is empty!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
		log.Infof("Set absolute path:%s, memory.limit_in_bytes path:%s\n", absolutePath, path.Join(absolutePath, "memory.limit_in_bytes"))
		if err := ioutil.WriteFile(path.Join(absolutePath, "memory.limit_in_bytes"), []byte(content), 0644); err != nil {
			log.Errorln("write content:" + content + "error!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
	}
	return nil
}

/**
将当前进程写入子的 cgroup mydocker 下的 tasks 文件中
*/
func (s *MemorySubsystem) Apply(pid string) error {
	absolutePath := ""
	if absolutePath = FindAbsolutePath(s.Name()); absolutePath == "" {
		log.Errorln("absolutePath is empty!")
		return fmt.Errorf("absolutePath is empty!\n")
	}
	log.Infof("Apply absolute path:%s, task path:%s\n", absolutePath, path.Join(absolutePath, "tasks"))
	if err := ioutil.WriteFile(path.Join(absolutePath, "tasks"), []byte(pid), 0644); err != nil {
		log.Errorln("write pid:" + pid + "error!")
		return fmt.Errorf("write pid:%s error!\n", pid)
	} else {
		log.Errorln("err:", err)
	}

	return nil
}

/**资源删除
删除对应的文件夹
*/
func (s *MemorySubsystem) Remove() error {
	absolutePath := ""
	if absolutePath = FindAbsolutePath(s.Name()); absolutePath == "" {
		log.Errorln("absolutePath is empty!")
		return fmt.Errorf("absolutePath is empty!\n")
	}
	if err := os.RemoveAll(absolutePath); err != nil {
		log.Errorln("remove absolute path error:", err)
		return fmt.Errorf("remove absolute path error:%v\n", err)
	}

	return nil
}
