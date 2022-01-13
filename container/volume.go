package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// NewWorkSpace 创建容器文件系统，进一步隔离容器和镜像，
// 实现容器中的操作不影响镜像；
func NewWorkSpace(volume, imageName, containerName string) {
	makeMntAndWriteLayer()
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs, containerName)
			logrus.Infof("New workspace volume urls %q", volumeURLs)
		} else {
			logrus.Infof("Volume parameter input is not correct.")
		}
	}
}

func CreateReadOnlyLayer(imageName string) error {
	unTarFoldURL := RootUrl + "/" + imageName + "/"
	imageUrl := RootUrl + "/" + imageName + ".tar"
	exist, err := PathExists(unTarFoldURL)
	if err != nil {
		logrus.Infof("Fail to judge whether dir %s exist: %v", unTarFoldURL, err)
	}
	if !exist {
		if err = os.Mkdir(unTarFoldURL, 0622); err != nil {
			logrus.Errorf("Mkdir dir %s error: %v", unTarFoldURL, err)
			return err
		}
		if _, err = exec.Command("tar", "-xvf", imageUrl, "-C", unTarFoldURL).CombinedOutput(); err != nil {
			logrus.Errorf("Untar dir %s error: %v", unTarFoldURL, err)
			return err
		}
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.Mkdir(writeURL, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("Mkdir dir %s error : %v", writeURL, err)
	}
}

func CreateMountPoint(containerName, imageName string) error {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.Mkdir(mntUrl, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("Mkdir dir %s error : %v", mntUrl, err)
	}
	tmpWriterLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	tmpImageLocation := RootUrl + "/" + imageName
	mntURL := fmt.Sprintf(MntUrl, containerName)
	dirs := "dirs=" + tmpWriterLayer + ":" + tmpImageLocation
	// mount -t aufs -o dirs=/root/writeLayer:/root/busybox none /root/mnt/
	// 使用 aufs 技术做到读写层分离，不影响实际镜像；
	// 参考链接 https://cloud.tencent.com/developer/article/1518056
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	logrus.Infof("cmd is %v", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Mount error : %v", err)
		return err
	}
	return nil
}

func MountVolume(volumeURLs []string, containerName string) error {
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil && os.IsNotExist(err) {
		logrus.Infof("Mkdir parent dir %s error : %v", parentUrl, err)
	}
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil && os.IsNotExist(err) {
		logrus.Infof("Mkdir container dir %s error : %v", containerVolumeURL, err)
	}
	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeURL).CombinedOutput()
	if err != nil {
		logrus.Errorf("Mount volume failed : %v", err)
		return err
	}
	return nil
}

func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		volumeURLs := strings.Split(volume, ":")
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(volumeURLs, containerName)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Errorf("Remove dir %s error : %v", writeURL, err)
	}
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	if _, err := exec.Command("umount", mntURL).CombinedOutput(); err != nil {
		logrus.Errorf("Umount error :  %v", err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("Remove mnt dir %s error : %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteMountPointWithVolume(volumeURLs []string, containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		logrus.Errorf("Umount volume failed : %v", err)
	}

	if _, err := exec.Command("umount", mntURL).CombinedOutput(); err != nil {
		logrus.Errorf("Umount mount point failed : %v", err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("Remove mount point dir failed : %v", err)
		return err
	}
	return nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		logrus.Errorf("Path: %s not exist", path)
		return false, err
	}
	return false, err
}

func makeMntAndWriteLayer() {
	mntDir, writeLayerDir := "/root/mnt", "/root/writeLayer"
	_, err := os.Stat(mntDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir("/root/mnt", 0622); err != nil {
			logrus.Errorf("Make mnt dir error : %v", err)
		}
	}
	_, err = os.Stat(writeLayerDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir("/root/writeLayer", 0622); err != nil {
			logrus.Errorf("Make writelayer error : %v", err)
		}
	}
}
