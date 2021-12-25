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
	// CLONE_NEWUTS is for new UTS namespace
	// CLONE_NEWIPC is for new IPC namespace
	// CLONE_NEWPID is for new PID namespace
	// CLONE_NEWNS is for new Mount Namespace,
	// please take care of use CLONE_NEWNS, because
	// you need to mount target director so that it
	// can work correctly.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
