package cgroups

import "mydocker/cgroups/subsystems"

type CgroupManger struct {
	// 限制的值
	Resource *subsystems.ResourceConfig
	// 用于在当前容器中标识有哪些 subsystem 需要做限制
	SubsystemsIns []subsystems.Subsystem
	// 用于启动多个 container 的 containerId。比如：memory 则是 /sys/fs/cgroup/memory/mydocker/[containerId]
	// path string
}

func (c *CgroupManger) Set() {
	for _, sub := range c.SubsystemsIns {
		sub.Set(c.Resource)
	}
}

func (c *CgroupManger) Apply(pid string) {
	for _, sub := range c.SubsystemsIns {
		sub.Apply(pid)
	}
}

func (c *CgroupManger) Destroy() {
	for _, sub := range c.SubsystemsIns {
		sub.Remove()
	}
}
