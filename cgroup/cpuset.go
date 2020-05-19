package cgroup

import (
	"io/ioutil"
	"kzdocker/log"
	"os"
	"path"
	"strconv"
)

type CPUSetSubsystem struct {
	path string
	res  *ResourceConfig
}

func NewCPUSetSubsystem(res *ResourceConfig) *CPUSetSubsystem {
	s := &CPUSetSubsystem{}
	s.path = findSubsystemMountpoint(s.Name())
	s.res = res
	return s
}

func (t *CPUSetSubsystem) Name() string {
	return "cpuset"
}

// Set 设置某个cgroup在这个Subsystem中的资源限制
func (t *CPUSetSubsystem) Set(cgroupPath string) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, true)
	if err != nil {
		return err
	}
	if t.res.CPUSet != "" {
		err = ioutil.WriteFile(path.Join(cpath, "cpuset.cpus"), []byte(t.res.CPUSet), 0644)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (t *CPUSetSubsystem) Apply(cgroupPath string, pid int) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, false)
	if err != nil {
		return err
	}
	// 这里需要在cpuset.mems中添加个0，不然回报错，详情看https://github.com/opencontainers/runc/issues/133
	if t.res.CPUSet != `` {
		err = ioutil.WriteFile(path.Join(cpath, "cpuset.mems"), []byte(`0`), 0644)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		err = ioutil.WriteFile(path.Join(cpath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (t *CPUSetSubsystem) Remove(cgroupPath string) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, false)
	if err != nil {
		return err
	}
	err = os.RemoveAll(cpath)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
