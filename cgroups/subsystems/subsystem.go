package subsystems

/*
用于传递资源限制配置的结构体，包含内存限制，CPU时间片权重，CPU核心数
*/
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

/*
Subsystem 接口，每个 Subsystem 可以实现下面的 4 个接口
这里将 cgroup 抽象成了 path，原因是 cgroup 在 hierarchy 的路径，便是虚拟文件系统中的虚拟路径
*/
type Subsystem interface {
	// 返回 Subsystem 的名字，比如 cpu memory。
	Name() string
	// 设置某个 cgroup 在这个 Subsystem 中的资源限制
	Set(path string, res *ResourceConfig) error
	// 将进程添加到某个 cgroup 中
	Apply(path string, pid int) error
	// 移除某个 cgroup
	Remove(path string) error
}

// 将不同的 subsystem 初始化实例创建资源限制处理链数组
var (
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
