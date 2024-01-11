package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	CgroupRoot = "/sys/fs/cgroup"
)

func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}
	return ""
}

func GetCgroupPath(cgroupPath string, autoCreate bool) (string, error) {
	// cgroupPath should be `/sys/fs/cgroup/{$containerName}`
	if _, err := os.Stat(path.Join(CgroupRoot, cgroupPath)); err == nil ||
		(autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(CgroupRoot, cgroupPath), 0755); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(CgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}

// ApplyCGroupProcess applies the specified pid to the cgroup specified by cGroupPath.
// It first gets the actual cgroup path using GetCgroupPath function.
// If an error occurs, it returns an error message indicating the failure.
// Then it writes the pid to the "cgroup.procs" file under the cgroupPath using the WriteFile function from the os package.
// If an error occurs, it returns an error message indicating the failure.
// If the operation is successful, it returns nil.
func ApplyCGroupProcess(cGroupPath string, pid int) error {
	cgroupPath, err := GetCgroupPath(cGroupPath, true)
	if err != nil {
		return fmt.Errorf("error getting cgroup path: %v", err)
	}
	err = os.WriteFile(path.Join(cgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("error applying cgroup process: %v", err)
	}
	return nil
}

// RemoveCGroup removes the cgroup specified by cGroupPath.
//
// It first gets the actual cgroup path using GetCgroupPath function.
// If an error occurs, it returns an error message indicating the failure.
//
// Then it removes the directory corresponding to the cgroup using RemoveAll function from the os package.
// If an error occurs, it returns an error message indicating the failure.
//
// If the operation is successful, it returns nil.
func RemoveCGroup(cGroupPath string) error {
	cgroupPath, err := GetCgroupPath(cGroupPath, true)
	if err != nil {
		return fmt.Errorf("error getting cgroup path: %v", err)
	}
	err = os.RemoveAll(cgroupPath)
	if err != nil {
		return fmt.Errorf("error removing cgroup: %v", err)
	}
	return nil
}
