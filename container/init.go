package container

import (
	"fmt"
	"os"
	"syscall"
)

// RunContainerInitProcess 新进程创建出来之后对新进程进行初始化
func RunContainerInitProcess(command string, args []string) error {
	fmt.Printf("command %s\n", command)

	syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), "")

	/*
		为什么要在容器中挂载/proc呢， 主要原因是因为ps、top等命令依赖于/proc目录。
		当隔离PID的时候，ps、top等命令还是未隔离的时候一样输出。 为了让隔离空间ps、top等命令只输出当前隔离空间的进程信息。需要单独挂载/proc目录。
	*/
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
