package subsystems

/**
这里定义了所有 subsystem 需要实现的接口 Subsystem 和资源变量 ResourceConfig
*/

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

type Subsystem interface {
	Name() string
	Set(res *ResourceConfig) error
	Apply(pid string) error
	Remove() error
}

const (
	ResourceName = "mydocker"
)
