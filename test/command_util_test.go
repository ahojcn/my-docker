package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"testing"
	"text/tabwriter"
	"time"
)

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}

var (
	RUNNING      string = "running"
	STOP         string = "STOP"
	EXIT         string = "exited"
	CONTAINS     string = "/var/run/mydocker"
	INFOLOCATION string = "/var/run/mydocker/%s"
	CONFIGNAME   string = "config.json"
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
	str := time.Now().UnixNano()
	containerId := fmt.Sprintf("%d%d", str, int(math.Abs(float64(rand.Intn(10)))))
	log.Infoln("containerId: %s\n", containerId)
	return containerId
}

// 保存容器 metadata，保存到 INFOLOCATION/uuid/config.json下
func RecordContainerInfo(pid, name, id, command string) error {
	containerInfo := &ContainerInfo{
		Pid:        pid,
		Id:         id,
		Name:       name,
		Command:    command,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:     RUNNING,
	}
	jsonInfo, _ := json.Marshal(containerInfo)
	log.Infoln("jsonInfo:", string(jsonInfo))
	location := fmt.Sprintf(INFOLOCATION, name)
	file := location + "/" + CONFIGNAME
	if err := os.MkdirAll(location, 0622); err != nil {
		return fmt.Errorf("create %s error: %v\n", location, err)
	}
	if err := ioutil.WriteFile(file, []byte(jsonInfo), 0622); err != nil {
		return fmt.Errorf("write %s to %s error: %v\n", jsonInfo, file, err)
	}
	return nil
}

// 获取容器 metadata
func GetContainerInfo(name string) (*ContainerInfo, error) {
	location := fmt.Sprintf(INFOLOCATION, name)
	file := location + "/" + CONFIGNAME
	containerInfo := &ContainerInfo{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read data %s error: %s\n", data, err)
	}
	json.Unmarshal(data, containerInfo)
	return containerInfo, nil
}

// 获取所有容器 metadata
func ShowAllContainers() {
	files, err := ioutil.ReadDir(CONTAINS)
	if err != nil {
		log.Errorln("read Dir error :", err)
		return
	}
	var containers []*ContainerInfo
	for _, file := range files {
		container, err := GetContainerInfo(file.Name())
		if err != nil {
			log.Errorln("error:", err)
			continue
		}
		containers = append(containers, container)
	}
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
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
	location := fmt.Sprintf(INFOLOCATION, name)
	if err := os.RemoveAll(location); err != nil {
		return fmt.Errorf("remove all %s error: %v\n", location, err)
	}
	return nil
}

func Test006(t *testing.T) {
	uuid := ContainerUUID()
	// 将 uuid 写入一个文件中，方便测试用一个 uuid 测试增删改查方法
	writeUUID(uuid)
}

func Test007(t *testing.T) {
	uuid := readUUID()
	if err := RecordContainerInfo(uuid, uuid, uuid, "/bin/top"); err != nil {
		log.Errorln("RecordContainerInfo error:", err)
	} else {
		log.Infoln("write successfully!")
	}
}

func Test008(t *testing.T) {
	uuid := readUUID()
	containerInfo, _ := GetContainerInfo(uuid)
	if containerInfo != nil {
		log.Infof("Pid:%s, Id:%s, Name:%s, Command:%s, CreateTime:%s, Status:%s\n",
			containerInfo.Pid, containerInfo.Id, containerInfo.Name, containerInfo.Command, containerInfo.CreateTime, containerInfo.Status)
	}
}

func Test009(t *testing.T) {
	ShowAllContainers()
}

func Test010(t *testing.T) {
	uuid := readUUID()
	if err := DeleteContainerInfo(uuid); err != nil {
		log.Errorln("DeleteContainerInfo error:", err)
	}
}
