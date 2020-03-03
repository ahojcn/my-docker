package command

type ContainerInfo struct {
	Pid        string `json:"pid"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Command    string `json:"command"`
	CreateTime string `json:"createTime"`
	Status     string `json:"status"`
}

var (
	RUNNING       string = "running"
	STOP          string = "STOP"
	EXIT          string = "exited"
	CONTAINS      string = "/var/run/mydocker"
	INFOLOCATION  string = "/var/run/mydocker/%s"
	CONFIGNAME    string = "config.json"
	CONTAINERLOGS string = "container.log"
)
