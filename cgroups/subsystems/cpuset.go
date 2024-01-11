package subsystems

import (
	"fmt"
	"os"
	"path"
)

type CpusetSubSystem struct{}

// Name 返回 subsystem 的名字 memory
func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}

// Set 设置 memory-cgroup 在这个 subsystem 中的资源限制
func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"),
				[]byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set cgroup cpuset fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}
