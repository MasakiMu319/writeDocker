package cgroups

import (
	"github.com/sirupsen/logrus"
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

func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsytemsIns {
		subSysIns.Apply(c.cgroupPath, pid)
	}
	return nil
}

func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsytemsIns {
		subSysIns.Set(c.cgroupPath, res)
	}
	return nil
}

func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsytemsIns {
		if err := subSysIns.Remove(c.cgroupPath); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
