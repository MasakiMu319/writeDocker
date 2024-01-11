package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
	"writeDocker/cgroups"
	"writeDocker/cgroups/subsystems"
	"writeDocker/container"
	"writeDocker/network"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume, containerName, imageName string,
	envSlice []string, nw string, portmapping []string) {
	containerId := randStringBytes(10)
	if containerName == "" {
		containerName = containerId
	}

	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName, envSlice)

	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Errorf("parent start error %v", err)
	}

	// generate containerName.
	containerName, err := recordContainerInfo(parent.Process.Pid, comArray, containerName, containerId, volume)
	if err != nil {
		logrus.Errorf("Record container info error : %v", err)
		return
	}

	cgroupManager := cgroups.NewCgroupManager(containerName)
	defer func(cGroupPath string) {
		err = subsystems.RemoveCGroup(cGroupPath)
		if err != nil {
			logrus.Errorf("Error removing cgroup: %v", err)
		}
	}(containerName)

	cgroupManager.Set(res)

	err = subsystems.ApplyCGroupProcess(containerName, parent.Process.Pid)
	if err != nil {
		logrus.Errorf("Error applying cgroup process: %v", err)
		return
	}

	if nw != "" {
		network.Init()
		containerInfo := &container.ContainerInfo{
			Id:          containerId,
			Pid:         strconv.Itoa(parent.Process.Pid),
			Name:        containerName,
			PortMapping: portmapping,
		}
		if err = network.Connect(nw, containerInfo); err != nil {
			logrus.Errorf("Error Connect Network %v", err)
			return
		}
	}

	sendInitCommand(comArray, writePipe)
	if tty {
		parent.Wait()
		container.DeleteWorkSpace(volume, containerName)
		deleteContainerInfo(containerName)
	}
	//parent.Wait()
}

// sendInitCommand send user's params value when init container.
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName, id, volume string) (string, error) {
	createTime := time.Now().Format("2006-01-02 15:04:05")
	command := strings.Join(commandArray, "")
	containerInfo := &container.ContainerInfo{
		Id:          id,
		Pid:         strconv.Itoa(containerPID),
		Command:     command,
		CreatedTime: createTime,
		Status:      container.RUNNING,
		Name:        containerName,
		Volume:      volume,
	}

	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Errorf("Record container info error : %v", err)
		return "", err
	}
	jsonStr := string(jsonBytes)

	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err = os.MkdirAll(dirUrl, 0622); err != nil {
		logrus.Errorf("Mkdir %s error : %v", dirUrl, err)
		return "", err
	}
	fileName := dirUrl + "/" + container.ConfigName
	file, err := os.Create(fileName)
	defer file.Close()
	if err != nil {
		logrus.Errorf("Create file %s error : %v", fileName, err)
		return "", err
	}
	if _, err = file.WriteString(jsonStr); err != nil {
		logrus.Errorf("File write string error %v", err)
		return "", err
	}
	return containerName, nil
}

func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func deleteContainerInfo(containerId string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerId)
	if err := os.RemoveAll(dirURL); err != nil {
		logrus.Errorf("Remove dir %s error : %v", dirURL, err)
	}
}
