package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// RunContainerInitProcess 容器初始化时发出系统调用先挂载一个独立的文件系统，
// 然后再挂载 /proc 目录，执行用户的 command；
//func RunContainerInitProcess(command string, args []string) error {
//	logrus.Infof("command %s", command)
//	// 先让新的 mount namespace 独立，否则退出运行后就会使得主机上的 /proc 需要重新 mount
//	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
//	// MS_NOEXEC 在本文件系统中不允许运行其他程序；
//	// MS_NOSUID 在本文件系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID；
//	// MS_NODEV 默认参数
//	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
//	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
//	argv := []string{command}
//	// syscall.Exec 执行 command 对应的程序，然后覆盖掉初始化程序中的上下文；
//	// 以此来实现容器的第一个进程是用户指定的程序；
//	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
//		logrus.Errorf(err.Error())
//	}
//	return nil
//}

// RunContainerInitProcess 容器初始化时发出系统调用先挂载一个独立的文件系统，
// 然后再挂载 /proc 目录，执行用户的 command；
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}
	// 先让新的 mount namespace 独立，否则退出运行后就会使得主机上的 /proc 需要重新 mount

	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// MS_NOEXEC 在本文件系统中不允许运行其他程序；
	// MS_NOSUID 在本文件系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID；
	// MS_NODEV 默认参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	// 调用 exec.LookPath 可以在系统的 PATH 中寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
