package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
	"github.com/Parad0xpl/git-remote-vesh/v2/external"
	"github.com/Parad0xpl/git-remote-vesh/v2/helper"
	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// ProcHandler represent process that can be started and stoped asynchronously
type ProcHandler interface {
	Start() error
	Stop() error
}

// MountSSHFS try to mount SSHFS resources with data from config. If VeraCrypt
// file is already accessible it is skipped.
func MountSSHFS(config config.VeshConfig) (ProcHandler, error) {
	log.Println("---[SSHFS Mount]---")

	if utils.IsDebug() {
		log.Println("Theoretical vault path:", config.VeraCryptVaultPath)
	}

	if _, err := os.Stat(config.VeraCryptVaultPath); err == nil {
		log.Println("Vault already found - skipping sshfs mount")
		return nil, nil
	}
	if _, err := os.Stat(config.SSHMountPath); err == nil {
		log.Println("SSH already mounted - skipping")
		return nil, nil
	}
	handle := external.CreateSSHFS(config.ExtractSSHFSParams())
	err := handle.Start()
	if err != nil {
		return nil, fmt.Errorf("can't mount sshfs: %e", err)
	}

	return handle, nil
}

// MountVeraCrypt try to mount VeraCrypt container specified in the config. If
// it is mounted already it is skipped.
func MountVeraCrypt(config config.VeshConfig) (ProcHandler, error) {
	log.Println("---[VeraCrypt Mount]---")

	if utils.IsDebug() {
		log.Println("Theoretical repo path:", config.RepoPath)
	}

	parent_dir := filepath.Dir(config.RepoPath)
	if _, err := os.Stat(parent_dir); err == nil {
		if _, err := os.Stat(config.RepoPath); err == nil {
			log.Println("Repo already found - skipping VeraCrypt mount")
		} else {
			log.Println("No repo but vera crypt folder found - creating bare repository")
			utils.CreateBareRepo(config.RepoPath)
		}
		return nil, nil
	}

	if !filepath.IsAbs(config.VeraCryptVaultPath) {
		prefix := config.SSHMountPath
		if strings.HasSuffix(prefix, ":") {
			prefix = prefix + "\\"
		}
		p := filepath.Join(prefix, config.VeraCryptVaultPath)

		config.VeraCryptVaultPath = p
	}

	if utils.IsDebug() {
		log.Println("VeraCrypt path:", config.VeraCryptVaultPath)
	}
	handle := external.CreateVeraCrypt(config.ExtractVeraCryptParams())
	err := handle.Start()
	if err != nil {
		return nil, err
	}

	return handle, nil
}

// DismountSSHFS dismount SSHFS if proc handler is available
func DismountSSHFS(handle ProcHandler) {
	log.Println("---[SSHFS Dismount]---")
	if handle != nil {
		handle.Stop()
	}
}

// DismountVeraCrypt dismount VeraCrypt if proc handler is available
func DismountVeraCrypt(handle ProcHandler) {
	log.Println("---[VeraCrypt Dismount]---")
	if handle != nil {
		handle.Stop()
	}
}

// Main realise base logic of application
func Main() error {
	if len(os.Args) != 3 {
		log.Printf("Usage: %s <remote name> <remote address>\n", os.Args[0])
		return nil
	}

	config, err := config.GetConfig()
	if err != nil {
		return err
	}

	sshHandle, err := MountSSHFS(config)
	if err != nil {
		return err
	}
	defer DismountSSHFS(sshHandle)
	veraHandle, err := MountVeraCrypt(config)
	if err != nil {
		return err
	}
	defer DismountVeraCrypt(veraHandle)

	log.Println("---[Begin Main]---")
	err = helper.Loop(config)
	if err != nil {
		return err
	}
	log.Println("---[End Main]---")

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
