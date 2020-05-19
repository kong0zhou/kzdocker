package container

import (
	"kzdocker/base"
	"kzdocker/log"
	"kzdocker/utils"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// NewParentProcess 创建出一个新进程，并添加NameSpace
func NewParentProcess(tty bool, command []string) (*exec.Cmd, error) {
	// args := []string{"init", command}
	// 调用自身
	cmd := exec.Command("/proc/self/exe", command...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	mntURL := filepath.Join(base.BasePath, `root/mnt`)
	rootURL := filepath.Join(base.BasePath, `root`)
	err := NewWorkSpace(rootURL, mntURL)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

//NewWorkSpace Create a AUFS filesystem as container root workspace
func NewWorkSpace(rootURL string, mntURL string) (err error) {
	err = CreateReadOnlyLayer(rootURL)
	if err != nil {
		return err
	}
	err = CreateWriteLayer(rootURL)
	if err != nil {
		return err
	}
	err = CreateMountPoint(rootURL, mntURL)
	if err != nil {
		return err
	}
	return nil
}

// CreateReadOnlyLayer 创建只读层
func CreateReadOnlyLayer(rootURL string) (err error) {
	alpineURL := filepath.Join(rootURL, "alpine")
	alpineTarURL := filepath.Join(rootURL, "alpine.tar")
	exist := utils.IsPathExist(alpineURL)
	if exist {
		return nil
	}
	if err = os.Mkdir(alpineURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", alpineURL, err)
		return err
	}
	if _, err = exec.Command("tar", "-xvf", alpineTarURL, "-C", alpineURL).CombinedOutput(); err != nil {
		log.Errorf("Untar dir %s error %v", alpineURL, err)
		return err
	}
	return nil
}

// CreateWriteLayer 创建可写层
func CreateWriteLayer(rootURL string) (err error) {
	// writeURL := rootURL + "writeLayer/"
	writeURL := filepath.Join(rootURL, "writeLayer")
	if err = os.Mkdir(writeURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL, err)
		return err
	}
	return nil
}

// CreateMountPoint !!!!创建挂载点
func CreateMountPoint(rootURL string, mntURL string) (err error) {
	if err = os.Mkdir(mntURL, 0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntURL, err)
		return err
	}
	dirs := "dirs=" + filepath.Join(rootURL, "writeLayer") + ":" + filepath.Join(rootURL, "alpine")
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	return nil
}

//DeleteWorkSpace Delete the AUFS filesystem while container exit
func DeleteWorkSpace(rootURL string, mntURL string) (err error) {
	err = DeleteMountPoint(mntURL)
	if err != nil {
		return err
	}
	DeleteWriteLayer(rootURL)
	if err != nil {
		return err
	}
	return nil
}

// DeleteMountPoint 删除挂载点
func DeleteMountPoint(mntURL string) (err error) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		log.Errorf("%v", err)
		return err
	}
	if err = os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

// DeleteWriteLayer 删除可写层
func DeleteWriteLayer(rootURL string) {
	// writeURL := rootURL + "writeLayer/"
	writeURL := filepath.Join(rootURL, "writeLayer")
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeURL, err)
	}
}
