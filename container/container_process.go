package container

import (
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 创建出一个新进程，并添加NameSpace
func NewParentProcess(tty bool, command []string) *exec.Cmd {
	// args := []string{"init", command}
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
	return cmd
}
