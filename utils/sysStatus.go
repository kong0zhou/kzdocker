package utils

import (
	"fmt"
	"runtime"
	"time"
)

// SysStatus 系统状态
type SysStatus struct {
	startTime time.Time //开始时间

	Uptime       string // 服务运行时间
	NumGoroutine int    // 当前 Goroutines 数量

	// General statistics.(常规统计信息)
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	Lookups      uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees

	// Main allocation heap statistics.
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects

	// Low-level fixed-size structure allocator statistics.
	//	Inuse is bytes used now.
	//	Sys is bytes obtained from system.
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations

	// Garbage collector statistics.
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string // GC 暂停时间总量
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32 // GC 执行次数
}

// NewSysStatus 新建一个系统状态监视器
func NewSysStatus() *SysStatus {
	return &SysStatus{
		startTime: time.Now(),
	}
}

// UpdateSystemStatus 更新系统状态
func (s *SysStatus) UpdateSystemStatus() {
	s.Uptime = time.Since(s.startTime).String()

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	s.NumGoroutine = runtime.NumGoroutine()

	s.MemAllocated = FileSize(int64(m.Alloc))
	s.MemTotal = FileSize(int64(m.TotalAlloc))
	s.MemSys = FileSize(int64(m.Sys))
	s.Lookups = m.Lookups
	s.MemMallocs = m.Mallocs
	s.MemFrees = m.Frees

	s.HeapAlloc = FileSize(int64(m.HeapAlloc))
	s.HeapSys = FileSize(int64(m.HeapSys))
	s.HeapIdle = FileSize(int64(m.HeapIdle))
	s.HeapInuse = FileSize(int64(m.HeapInuse))
	s.HeapReleased = FileSize(int64(m.HeapReleased))
	s.HeapObjects = m.HeapObjects

	s.StackInuse = FileSize(int64(m.StackInuse))
	s.StackSys = FileSize(int64(m.StackSys))
	s.MSpanInuse = FileSize(int64(m.MSpanInuse))
	s.MSpanSys = FileSize(int64(m.MSpanSys))
	s.MCacheInuse = FileSize(int64(m.MCacheInuse))
	s.MCacheSys = FileSize(int64(m.MCacheSys))
	s.BuckHashSys = FileSize(int64(m.BuckHashSys))
	s.GCSys = FileSize(int64(m.GCSys))
	s.OtherSys = FileSize(int64(m.OtherSys))

	s.NextGC = FileSize(int64(m.NextGC))
	s.LastGC = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	s.PauseTotalNs = fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000)
	s.PauseNs = fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000)
	s.NumGC = m.NumGC
}
