package subsystems

import (
	"fmt"
	"os"
	"path"
)

type MemorySubSystem struct{}

// Name 返回 subsystem 的名字 memory
func (s *MemorySubSystem) Name() string {
	return "memory"
}

// Set 设置 memory-cgroup 在这个 subsystem 中的资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			if err := os.WriteFile(path.Join(subsysCgroupPath, "memory.max"),
				[]byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}
