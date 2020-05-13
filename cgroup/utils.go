package cgroup

import (
	"bufio"
	"fmt"
	"kzdocker/log"
	"kzdocker/utils"
	"os"
	"path"
	"strings"
)

// findSubsystemMountpoint 找到subsystems的挂载目录
func findSubsystemMountpoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

func getCgroupPath(subsystemPath string, cgroupPath string, autoCreate bool) (cpath string, err error) {
	// Join joins any number of path elements into a single path, adding a separating slash if necessary. The result is Cleaned; in particular, all empty strings are ignored.
	// if _, err := os.Stat(path.Join(subsystemPath, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
	// 	if os.IsNotExist(err) {
	// 		if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
	// 		} else {
	// 			return "", fmt.Errorf("error create cgroup %v", err)
	// 		}
	// 	}
	// 	return path.Join(cgroupRoot, cgroupPath), nil
	// } else {
	// 	return "", fmt.Errorf("cgroup path error %v", err)
	// }
	cpath = path.Join(subsystemPath, cgroupPath)
	if !utils.IsPathExist(cpath) {
		log.Info(`cpath is not exist`)
		if !autoCreate {
			err = fmt.Errorf(`path is not exist and not create`)
			log.Error(err.Error())
			return ``, err
		}
		err = os.Mkdir(cpath, 0755)
		if err != nil {
			log.Error(err.Error())
			return ``, err
		}
	}
	return cpath, nil
}
