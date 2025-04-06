package config

import (
	"fmt"
	"os"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// IsVeraCryptVaultAvailable checks if the VeraCrypt file is in place. It verifies the existence,
// file type, and signature of the file.
func (config *VeshConfig) IsVeraCryptVaultAvailable() error {
	stat, err := os.Stat(config.VeraCryptVaultPath)
	if err != nil {
		return fmt.Errorf("can't stat VeraCrypt vault: %w", err)
	}

	if !stat.Mode().IsRegular() {
		return fmt.Errorf("VeraCrypt vault is not a regular file")
	}

	return nil
}

// CheckSSHFS performs a safe check to ensure everything is in place after the SSHFS mount:
// * SSHFS is properly mounted.
// * The VeraCrypt vault file exists and is a regular file.
func (config *VeshConfig) CheckSSHFS() error {
	if !config.IsSSHFSMounted() {
		return fmt.Errorf("SSHFS is not mounted properly")
	}

	if err := config.IsVeraCryptVaultAvailable(); err != nil {
		return fmt.Errorf("problem with VeraCrypt vault: %w", err)
	}

	return nil
}

// CheckVeraCrypt performs a safe check to ensure everything is properly mounted.
func (config *VeshConfig) CheckVeraCrypt() error {
	if !config.IsVeraCryptMounted() {
		return fmt.Errorf("VeraCrypt is not mounted properly")
	}

	stat, err := os.Stat(config.RepoPath)
	if err != nil {
		if os.IsNotExist(err) {
			utils.CreateBareRepo(config.RepoPath)
			return nil
		}

		return fmt.Errorf("problem with repo: %w", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("repo is not a directory")
	}

	return nil
}
