package cgroups

import (
	"writeDocker/cgroups/subsystems"
)

type CgroupManager struct {
	cgroupPath string
	res        *subsystems.ResourceConfig
}

func NewCgroupManager(cgroupPath string) *CgroupManager {
	return &CgroupManager{
		cgroupPath: cgroupPath,
	}
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsytemsIns {
		subSysIns.Set(c.cgroupPath, res)
	}
	return nil
}
