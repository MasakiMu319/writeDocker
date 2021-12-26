package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"syscall"
)

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)
	// 先让新的 mount namespace 独立，否则退出运行后就会使得主机上的 /proc 需要重新 mount
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// MS_NOEXEC 在本文件系统中不允许运行其他程序；
	// MS_NOSUID 在本文件系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID；
	// MS_NODEV 默认参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}
