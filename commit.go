package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"writeDocker/container"
)

func commitContainer(containerName, imageName string) {
	mntURL := fmt.Sprintf(container.MntUrl, containerName)
	mntURL += "/"
	imageTar := container.RootUrl + "/" + imageName + ".tar"
	if _, err := exec.Command("tar", "-czf",
		imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error : %v", mntURL, err)
	}
}
