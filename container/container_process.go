// TODO: remember add os match here.

package container

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
)

// NewParentProcess 每次从当前进程的运行环境中 fork 一个新的进程，
// 并使用 namespace 进行初始化；
func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	// 传入管道文件读取端的句柄；
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = "/root/busybox"
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	// 生成一个匿名管道，读写变量都是文件类型；
	// 与 Linux 系统管道的定义保持一致；
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
