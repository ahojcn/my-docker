package container

type ContainerInfo struct {
	Pid         string   `json:"pid"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Command     string   `json:"command"`
	CreateTime  string   `json:"createTime"`
	Status      string   `json:"status"`
	Volumes     []string `json:"volumes"`
	RootPath    string   `json:"rootPaths"`
	PortMapping []string `json:"portmapping"` //端口映射
}

var (
	RUNNING       string = "running"
	STOP          string = "stopped"
	EXIT          string = "exited"
	CONTAINS      string = "/var/run/mydocker"
	INFOLOCATION  string = "/var/run/mydocker/%s"
	CONFIGNAME    string = "config.json"
	CONTAINERLOGS string = "container.log"
)
