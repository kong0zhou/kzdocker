package cgroup

// ResourceConfig 用于传递资源限制配置的结构体
type ResourceConfig struct {
	MemoryLimit string //包含内存限制
	CPUShare    string //cpu时间片权重
	CPUSet      string //cpu核心数
}

//Subsystem 接口，每个Subsystem可以实现下面4个接口
// 这里将cgroup抽象成path,原因是cgroup在hiearchy的路径，便是虚拟系统中的虚拟路径
type Subsystem interface {
	// 返回subsystem的名字
	Name() string
	// 设置某个cgroup在这个Subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// 将进程添加到某个cgroup中
	Apply(path string, pid int) error
	// 移除某个cgroup
	Remove(path string) error
}
