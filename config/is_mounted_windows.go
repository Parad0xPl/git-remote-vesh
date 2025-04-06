//go:build windows

package config

import (
	"os"
	"path/filepath"
)

// IsVeraCryptMounted checks if VeraCrypt is already mounted.
func (config *VeshConfig) IsVeraCryptMounted() bool {
	parentDir := filepath.Dir(config.RepoPath)
	if _, err := os.Stat(parentDir); err == nil {
		return true
	}

	return false
}

// IsSSHFSMounted checks if SSHFS is already mounted.
func (config *VeshConfig) IsSSHFSMounted() bool {
	if _, err := os.Stat(config.SSHMountPath); err == nil {
		return true
	}

	return false
}
