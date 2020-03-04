package command

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"mydocker/container"
	"os"
	"text/tabwriter"
	"time"
)

func writeUUID(uuid string) {
	ioutil.WriteFile("uuid.txt", []byte(uuid), 0644)
}

func readUUID() string {
	data, _ := ioutil.ReadFile("uuid.txt")
	return string(data)
}

// 生成 uuid
func ContainerUUID() string {
	//str := time.Now().UnixNano()
	//containerId := fmt.Sprintf("%d%d", str, int(math.Abs(float64(rand.Intn(10)))))
	//log.Infoln("containerId:", containerId)
	//return containerId
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// 保存容器 metadata，保存到 INFOLOCATION/uuid/config.json下
func RecordContainerInfo(pid, name, id, command string, volumes []string, newRootPath string) error {
	containerInfo := &container.ContainerInfo{
		Pid:        pid,
		Id:         id,
		Name:       name,
		Command:    command,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     container.RUNNING,
		Volumes:    volumes,
		RootPath:   newRootPath,
	}
	jsonInfo, _ := json.Marshal(containerInfo)
	log.Infoln("jsonInfo:", string(jsonInfo))
	location := fmt.Sprintf(container.INFOLOCATION, name)
	file := location + "/" + container.CONFIGNAME
	if err := os.MkdirAll(location, 0622); err != nil {
		return fmt.Errorf("create %s error: %v\n", location, err)
	}
	if err := ioutil.WriteFile(file, []byte(jsonInfo), 0622); err != nil {
		return fmt.Errorf("write %s to %s error: %v\n", jsonInfo, file, err)
	}
	return nil
}

// 获取容器 metadata
func GetContainerInfo(name string) (*container.ContainerInfo, error) {
	location := fmt.Sprintf(container.INFOLOCATION, name)
	file := location + "/" + container.CONFIGNAME
	containerInfo := &container.ContainerInfo{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read data %s error: %s", data, err)
	}
	json.Unmarshal(data, containerInfo)
	return containerInfo, nil
}

// 获取所有容器 metadata
func ShowAllContainers() {
	files, err := ioutil.ReadDir(container.CONTAINS)
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")

	if err != nil {
		log.Errorln("ReadDir error :", err)
		return
	}
	var containers []*container.ContainerInfo
	for _, file := range files {
		container, err := GetContainerInfo(file.Name())
		if err != nil {
			log.Errorln("error:", err)
			continue
		}
		containers = append(containers, container)
	}
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreateTime)
	}
	if err := w.Flush(); err != nil {
		log.Errorln("flush error:", err)
	}
}

// 删除容器 metadata
func DeleteContainerInfo(name string) error {
	location := fmt.Sprintf(container.INFOLOCATION, name)
	if err := os.RemoveAll(location); err != nil {
		return fmt.Errorf("remove all %s error: %v\n", location, err)
	}
	return nil
}

// 创建 log file
func GetLogFile(containerName string) (*os.File, error) {
	path := fmt.Sprintf(container.INFOLOCATION, containerName)
	logFile := path + "/" + container.CONTAINERLOGS
	if err := os.MkdirAll(path, 0622); err != nil {
		return nil, fmt.Errorf("create %s error: %v\n", path, err)
	}
	if file, err := os.Create(logFile); err != nil {
		return nil, fmt.Errorf("os.Create %s error: %v\n", logFile, err)
	} else {
		return file, nil
	}
}

// 读取 log
func ReadLogs(containerName string) string {
	path := fmt.Sprintf(container.INFOLOCATION, containerName)
	logFile := path + "/" + container.CONTAINERLOGS
	data, _ := ioutil.ReadFile(logFile)
	return string(data)
}

// 更新容器 metadata
func UpdateContainerInfo(containerInfo *container.ContainerInfo) error {
	jsonInfo, _ := json.Marshal(containerInfo)
	//log.Printf("jsonInfo:%s\n", string(jsonInfo))
	location := fmt.Sprintf(container.INFOLOCATION, containerInfo.Name)
	file := location + "/" + container.CONFIGNAME
	if err := ioutil.WriteFile(file, []byte(jsonInfo), 0622); err != nil {
		return fmt.Errorf("write %s to %s error:%v\n", jsonInfo, file, err)
	}
	return nil
}

// 删除容器信息
func RemoveContainerInfo(containerInfo *container.ContainerInfo) error {
	location := fmt.Sprintf(container.INFOLOCATION, containerInfo.Name)
	if err := os.RemoveAll(location); err != nil {
		return fmt.Errorf("os.RemoveAll(%s) %v\n", location, err)
	}
	return nil
}
