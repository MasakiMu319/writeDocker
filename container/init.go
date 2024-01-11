//go:build linux
// +build linux

package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// RunContainerInitProcess 容器初始化时发出系统调用先挂载一个独立的文件系统，
// 然后再挂载 /proc 目录，执行用户的 command；
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	setUpMount()

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

func setUpMount() {
	root, err := os.Getwd()
	if err != nil {
		logrus.Errorf("Get current location error :%v", err)
		return
	}
	logrus.Infof("Current location is %s ", root)

	err = pivotRoot(root)
	if err != nil {
		logrus.Errorf("pivot root error: %v", err)
	}
	// 先让新的 mount namespace 独立，否则退出运行后就会使得主机上的 /proc 需要重新 mount
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// MS_NOEXEC 在本文件系统中不允许运行其他程序；
	// MS_NOSUID 在本文件系统中运行程序的时候，不允许 set-user-ID 或 set-group-ID；
	// MS_NODEV 默认参数
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

// pivotRoot is a function that changes the root directory of the current process to the specified directory.
// It performs the following steps:
//  1. Mounts a new root file system using the specified directory as both the source and target.
//     This is done to ensure that the old and new root are not on the same file system.
//  2. Creates a directory named ".pivot_root" inside the new root to store the old root.
//     If the directory already exists, it is not recreated.
//  3. Calls the PivotRoot system call to switch the root file system to the new root.
//     The old root is mounted at "root/.pivot_root".
//  4. Changes the current working directory to the new root.
//  5. Unmounts the old root from "root/.pivot_root".
//  6. Deletes the ".pivot_root" directory.
//
// The function takes the `root` string parameter, which specifies the new root directory path.
// It returns an error if any of the steps fail, or nil if successful.
func pivotRoot(root string) error {
	/** mount --bind ${root} ${root}
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("Mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	logrus.Infof("mkdir")
	// 需要注意，这里如果挂载点存放的文件夹已经存在 os.Mkdir 会返回 err
	// 所以，这里判断 err 是不是已经创建来判断是否需要返回 err
	if err := os.Mkdir(pivotDir, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("mkdir error: %v", err)
		return err
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	// debug：logrus.Infof("pivot_root: root: %s; pivotDir: %s", root, pivotDir)
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot_root %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}
	logrus.Infof("done pivot_root")
	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
