package subsystems

// ResourceConfig 传递资源配置信息的结构体， 这里
// 包含了内存限制，CPU 时间片权重，CPU 核心数目；
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// Subsystem 接口，将 cgroup 抽象成为 cgroupPath；
type Subsystem interface {
	// Name 返回对应 subsystem 的名字，比如 cpu，memory
	Name() string
	// Set 设置某个 cgroup 在这个 subsystem 中的资源限制
	Set(cgroupPath string, res *ResourceConfig) error
}

var (
	SubsytemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
