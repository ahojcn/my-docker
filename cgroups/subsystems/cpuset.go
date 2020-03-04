package subsystems

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

type CpusetSubSystem struct {
}

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}

/**
向 cpu.shares 写入限制 content
*/
func (s *CpusetSubSystem) Set(res *ResourceConfig) error {
	if res.CpuShare != "" {
		content := res.CpuShare
		absolutePath := ""
		if absolutePath = FindAbsolutePath(s.Name()); absolutePath == "" {
			log.Errorln("absolutePath is empty!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
		log.Infof("Set absolute path:%s, cpuset.cpus path:%s\n", absolutePath, path.Join(absolutePath, "memory.limit_in_bytes"))
		if err := ioutil.WriteFile(path.Join(absolutePath, "cpuset.cpus"), []byte(content), 0644); err != nil {
			log.Errorln("write content:" + content + "error!")
			return fmt.Errorf("absolutePath is empty!\n")
		}
	}
	return nil
}

/**
将当前进程写入子的 cgroup mydocker 下的 tasks 文件中
*/
func (s *CpusetSubSystem) Apply(pid string) error {
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
func (s *CpusetSubSystem) Remove() error {
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

