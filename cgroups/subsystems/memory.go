package subsystems

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"path"
)

/**
向 memory.limit_in_bytes 写入限制 content
*/
func Set(content string) error {
	absolutePath := ""
	if absolutePath = FindAbsolutePath("memory"); absolutePath == "" {
		log.Errorln("absolutePath is empty!")
		return fmt.Errorf("absolutePath is empty!\n")
	}
	log.Infof("Apply absolute path:%s, memory.limit_in_bytes path:%s\n", absolutePath, path.Join(absolutePath, "memory.limit_in_bytes"))
	if err := ioutil.WriteFile(path.Join(absolutePath, "memory.limit_in_bytes"), []byte(content), 0644); err != nil {
		log.Errorln("write content:" + content + "error!")
		return fmt.Errorf("absolutePath is empty!\n")
	}
	return nil
}

/**
将当前进程写入子的 cgroup mydocker 下的 tasks 文件中
*/
func Apply(pid string) error {
	absolutePath := ""
	if absolutePath = FindAbsolutePath("memory"); absolutePath == "" {
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
