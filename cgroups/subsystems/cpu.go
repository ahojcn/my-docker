package subsystems

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type CpuSubSystem struct {
}

func (s *CpuSubSystem) Name() string {
	return "cpu"
}

/**
向 cpu.shares 写入限制 content
*/
func (s *CpuSubSystem) Set(res *ResourceConfig) error {
	if res.CpuShare != "" {
		content := res.CpuShare
		absolutePath := ""
		if absolutePath = FindAbsolutePath(s.Name()); absolutePath == "" {
			log.Errorln("absolutePath is empty!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
		log.Infof("Set absolute path:%s, cpu.shares path:%s\n", absolutePath, path.Join(absolutePath, "memory.limit_in_bytes"))
		if err := ioutil.WriteFile(path.Join(absolutePath, "cpu.shares"), []byte(content), 0644); err != nil {
			log.Errorln("write content:" + content + "error!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
	}
	return nil
}

/**
将当前进程写入子的 cgroup mydocker 下的 tasks 文件中
*/
func (s *CpuSubSystem) Apply(pid string) error {
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
func (s *CpuSubSystem) Remove() error {
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
