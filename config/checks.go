package config

import (
	"fmt"
	"os"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// IsVeraCryptVaultAvailable check if VeraCrypt file is in place. Check for exsitance,
// filetype and signature.
func (config *VeshConfig) IsVeraCryptVaultAvailable() error {
	stat, err := os.Stat(config.VeraCryptVaultPath)
	if err != nil {
		return fmt.Errorf("can't stat vera crypt vault: %e", err)
	}

	if !stat.Mode().IsRegular() {
		return fmt.Errorf("VeraCrypt vault is not a regular file")
	}

	return nil
}

// CheckSSHFS do safe check if all things after sshfs mount are in place:
// * SSHFS is realy mounter
// * VeraCrypt Vault file exists and is a regular file
func (config *VeshConfig) CheckSSHFS() error {
	if !config.IsSSHFSMounted() {
		return fmt.Errorf("SSHFS is not mounted properly")
	}

	if err := config.IsVeraCryptVaultAvailable(); err != nil {
		return fmt.Errorf("problem with VeraCrypt vault: %e", err)
	}

	return nil
}

// CheckVeraCrypt do safe check if all is properly mounted.
func (config *VeshConfig) CheckVeraCrypt() error {
	if !config.IsVeraCryptMounted() {
		return fmt.Errorf("VeraCrypt is not mounted properly")
	}

	stat, err := os.Stat(config.RepoPath)
	if err != nil {
		if !os.IsExist(err) {
			utils.CreateBareRepo(config.RepoPath)
			return nil
		}

		return fmt.Errorf("problem with repo: %e", err)
	}

	if !stat.IsDir() {
		return fmt.Errorf("repo is not a directory")
	}

	return nil
}
