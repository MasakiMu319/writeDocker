package subsystems

import (
	"fmt"
	"os"
	"path"
)

type CpuSubSystem struct{}

// Name 返回 subsystem 的名字 cpu
func (s *CpuSubSystem) Name() string {
	return "cpu"
}

// Set 设置 cpu-cgroup 在这个 subsystem 中的资源限制
func (s *CpuSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(cgroupPath, true); err == nil {
		if res.CpuShare != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.max"),
				[]byte(res.CpuShare), 0644); err != nil {
				return fmt.Errorf("set cgroup cpu fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}
