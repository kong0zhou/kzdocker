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
}

func NewMemorySubsystem() *MemorySubsystem {
	s := &MemorySubsystem{}
	s.path = findSubsystemMountpoint(s.Name())
	return s
}

func (t *MemorySubsystem) Name() string {
	return "memory"
}

// Set 设置某个cgroup在这个Subsystem中的资源限制
func (t *MemorySubsystem) Set(cgroupPath string, res *ResourceConfig) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, true)
	if err != nil {
		return err
	}
	if res.MemoryLimit != "" {
		err = ioutil.WriteFile(path.Join(cpath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644)
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
