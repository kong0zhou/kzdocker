package cgroup

import "kzdocker/log"

type CGroupManager struct {
	// cgroup在hierarchy中的路径 相当于创建的cgroup目录相对于root cgroup目录的路径
	path string
	// 资源配置
	resource *ResourceConfig
	//
	subsystemsIns []Subsystem
}

func NewCGroupManager(path string, res *ResourceConfig) *CGroupManager {
	return &CGroupManager{
		path:     path,
		resource: res,
		subsystemsIns: []Subsystem{
			NewMemorySubsystem(res),
			NewCPUSubsystem(res),
			NewCPUSetSubsystem(res),
		},
	}
}

func (c *CGroupManager) Set() (err error) {
	for _, subSysIns := range c.subsystemsIns {
		err = subSysIns.Set(c.path)
		if err != nil {
			return err
		}
	}
	return nil
}

// Apply 将进程pid加入到这个cgroup中
func (c *CGroupManager) Apply(pid int) (err error) {
	for _, subSysIns := range c.subsystemsIns {
		err = subSysIns.Apply(c.path, pid)
		if err != nil {
			return err
		}
	}
	return nil
}

// Destroy 释放cgroup
func (c *CGroupManager) Destroy() (err error) {
	log.Info(`begin destory`)
	for _, subSysIns := range c.subsystemsIns {
		if err := subSysIns.Remove(c.path); err != nil {
			log.Warnf("remove cgroup fail %v", err)
			return err
		}
	}
	return nil
}
