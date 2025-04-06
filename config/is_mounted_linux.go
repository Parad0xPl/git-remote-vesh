//go:build linux

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
)

type mountInfoData struct {
	path string
}

// getMountinfoData retrieves mount information from /proc/self/mountinfo.
func getMountinfoData() ([]mountInfoData, error) {
	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return nil, fmt.Errorf("Can't read info about mounts: %e", err)
	}

	mountInfo := string(data)
	mountInfo = strings.TrimSpace(mountInfo)
	lines := strings.Split(mountInfo, "\n")
	mountData := make([]mountInfoData, len(lines))
	for index, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) > 5 {
			mountData[index] = mountInfoData{
				path: parts[4],
			}
		}
	}

	return mountData, nil
}

// IsVeraCryptMounted checks if VeraCrypt is already mounted.
func (config *VeshConfig) IsVeraCryptMounted() bool {
	mountData, err := getMountinfoData()
	if err != nil {
		return false
	}

	debug.Println("VeraCrypt mount path:", config.VeraCryptMountPath)
	for _, d := range mountData {
		if d.path == config.VeraCryptMountPath {
			return true
		}
	}
	return false
}

// IsSSHFSMounted checks if SSHFS is already mounted.
func (config *VeshConfig) IsSSHFSMounted() bool {
	mountData, err := getMountinfoData()
	if err != nil {
		return false
	}

	for _, d := range mountData {
		if d.path == config.SSHMountPath {
			return true
		}
	}
	return false
}
