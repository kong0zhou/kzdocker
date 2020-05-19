package cgroup

import (
	"io/ioutil"
	"kzdocker/log"
	"os"
	"path"
	"strconv"
)

type MemorySubsystem struct {
	path string
	res  *ResourceConfig
}

func NewMemorySubsystem(res *ResourceConfig) *MemorySubsystem {
	s := &MemorySubsystem{}
	s.path = findSubsystemMountpoint(s.Name())
	s.res = res
	return s
}

func (t *MemorySubsystem) Name() string {
	return "memory"
}

// Set 设置某个cgroup在这个Subsystem中的资源限制
func (t *MemorySubsystem) Set(cgroupPath string) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, true)
	if err != nil {
		return err
	}
	if t.res.MemoryLimit != "" {
		err = ioutil.WriteFile(path.Join(cpath, "memory.limit_in_bytes"), []byte(t.res.MemoryLimit), 0644)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (t *MemorySubsystem) Apply(cgroupPath string, pid int) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, false)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(cpath, "tasks"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (t *MemorySubsystem) Remove(cgroupPath string) (err error) {
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
