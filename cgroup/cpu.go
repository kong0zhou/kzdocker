package cgroup

import (
	"io/ioutil"
	"kzdocker/log"
	"os"
	"path"
	"strconv"
)

type CPUSubsystem struct {
	path string
}

func NewCPUSubsystem() *CPUSubsystem {
	s := &CPUSubsystem{}
	s.path = findSubsystemMountpoint(s.Name())
	return s
}

func (t *CPUSubsystem) Name() string {
	return "cpu"
}

// Set 设置某个cgroup在这个Subsystem中的资源限制
func (t *CPUSubsystem) Set(cgroupPath string, res *ResourceConfig) (err error) {
	cpath, err := getCgroupPath(t.path, cgroupPath, true)
	if err != nil {
		return err
	}
	if res.CPUShare != "" {
		err = ioutil.WriteFile(path.Join(cpath, "cpu.shares"), []byte(res.CPUShare), 0644)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}

func (t *CPUSubsystem) Apply(cgroupPath string, pid int) (err error) {
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

func (t *CPUSubsystem) Remove(cgroupPath string) (err error) {
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
