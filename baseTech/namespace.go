//go:build linux
// +build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	// CLONE_NEWUTS is for new UTS namespace;
	// CLONE_NEWIPC is for new IPC namespace;
	// CLONE_NEWPID is for new PID namespace;
	// CLONE_NEWNS is for new Mount namespace,
	// please take care of use CLONE_NEWNS, because
	// you need to mount target director so that it
	// can work correctly;
	// CLONE_NEWUSER is for new user group namespace;
	// CLONE_NEWNET is for new Network namespace;
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

// tips1: if you exec "mount -t proc proc /proc" in namespace,
// remember exec it again after you back to host.
