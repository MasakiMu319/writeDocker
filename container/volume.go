package container

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func NewWorkSpace(volume, imageName, containerName string) {
	makeWorkspace()
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

// CreateMountPoint creates the mount point for a container, using overlayfs to isolate the container from the image.
// The function takes the container name and image name as parameters.
// It creates the necessary directories for the mount point, including the lower, upper, and work directories.
// Then it constructs the overlayfs options with the image location as the lower directory, the writer layer as the upper directory,
// and the work directory as the work directory.
// Finally, it runs the "mount" command to mount the overlayfs with the constructed options to the mount point directory.
// If successful, it returns nil. Otherwise, it returns an error.
func CreateMountPoint(containerName, imageName string) error {
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.Mkdir(mntUrl, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("Mkdir dir %s error : %v", mntUrl, err)
	}
	tmpImageLocation := RootUrl + "/" + imageName
	tmpWriterLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.Mkdir(tmpWriterLayer, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("Mkdir dir %s error : %v", tmpWriterLayer, err)
	}
	tmpWorkDir := fmt.Sprintf(WorkDirUrl, containerName)
	if err := os.Mkdir(tmpWorkDir, 0777); err != nil && os.IsNotExist(err) {
		logrus.Errorf("Mkdir dir %s error : %v", tmpWorkDir, err)
	}
	dirs := "lowerdir=" + tmpImageLocation + ",upperdir=" + tmpWriterLayer + ",workdir=" + tmpWorkDir
	mntURL := fmt.Sprintf(MntUrl, containerName)
	// mount -t overlay -o lowerdir=/root/busybox,upperdir=/root/writeLayer,workdir=/root/workdir overlay /root/mnt/
	// Update date: 2024.01 Change aufs with overlay, the aufs has decreased .
	cmd := exec.Command("mount", "-t", "overlay", "-o", dirs, "overlay", mntURL)
	logrus.Infof("cmd is : %v", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("Mount error : %v", err)
		return err
	}
	return nil
}

// MountVolume creates a mount point for a volume in a container.
// The function takes a slice of volume URLs and the container name as parameters.
// It creates the necessary directories for the mount point, including the host volume directory and the container directory.
// Then it uses the "mount --bind" command to mount the host volume directory to the container directory.
// If successful, it returns nil. Otherwise, it returns an error.
func MountVolume(volumeURLs []string, containerName string) error {
	// host volume directory.
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil && os.IsNotExist(err) {
		logrus.Infof("Mkdir parent dir %s error : %v", parentUrl, err)
	}
	// container directory.
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerVolumeURL := mntURL + "/" + containerUrl
	if err := os.Mkdir(containerVolumeURL, 0777); err != nil && os.IsNotExist(err) {
		logrus.Infof("Mkdir container dir %s error : %v", containerVolumeURL, err)
	}
	// use mount --bind instead of overlay FS.
	cmd := exec.Command("mount", "--bind", parentUrl, containerVolumeURL)
	logrus.Infof("cmd is : %v", cmd)

	if err := cmd.Run(); err != nil {
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

// makeWorkspace is a function that creates workspace directories for a container.
// It creates three directories: mntDir, writeLayerDir, and workDir.
// If any of these directories do not exist, the function creates them with 0622 permission.
// The mntDir directory is created at "/root/mnt", the writeLayerDir directory is created at "/root/writeLayer",
// and the workDir directory is created at "/root/workDir".
//
// Example usage:
// makeWorkspace()
func makeWorkspace() {
	mntDir, writeLayerDir, workDir := "/root/mnt", "/root/writeLayer", "/root/workDir"
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
	_, err = os.Stat(workDir)
	if os.IsNotExist(err) {
		if err = os.Mkdir("/root/workDir", 0622); err != nil {
			logrus.Errorf("Make work dir error : %v", err)
		}
	}
}
