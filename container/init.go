package container

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// RunContainerInitProcess 新进程创建出来之后对新进程进行初始化
func RunContainerInitProcess(command []string, args []string) error {
	fmt.Printf("command %v\n", command)

	syscall.Mount("", "/", "", uintptr(syscall.MS_PRIVATE|syscall.MS_REC), "")
	/*

			为什么要在容器中挂载/proc呢， 主要原因是因为ps、top等命令依赖于/proc目录。
			当隔离PID的时候，ps、top等命令还是未隔离的时候一样输出。 为了让隔离空间ps、top等命令只输出当前隔离空间的进程信息。需要单独挂载/proc目录。

		defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
		syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	*/

	setUpMount()

	// SYS EXECVE 系统调用井不会去帮我们在PATH中寻找命令，所以要调用exec.LookPath
	path, err := exec.LookPath(command[0])
	if err != nil {
		fmt.Println(`exec.LookPath failed:`, err.Error())
		return err
	}
	fmt.Println(`path is:`, path)
	if err := syscall.Exec(path, command[0:], os.Environ()); err != nil {
		fmt.Println(`syscall.Exec error:`, err.Error())
		return err
	}
	return nil
}

/**
Init 挂载点
*/
func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("setUpMount() Get current location error " + err.Error())
		// return
	}
	fmt.Println("setUpMount() Current location is ", pwd)
	err = pivotRoot(filepath.Join(pwd, "root/mnt"))
	if err != nil {
		panic("setUpMount()  error " + err.Error())
		// return
	}

	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
