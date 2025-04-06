package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
	"github.com/Parad0xpl/git-remote-vesh/v2/external"
	"github.com/Parad0xpl/git-remote-vesh/v2/helper"
	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// ProcHandler represents a process that can be started and stopped asynchronously
type ProcHandler interface {
	Start() error
	Stop() error
}

// MountSSHFS attempts to mount SSHFS resources using data from the config. If the VeraCrypt
// file is already accessible, it is skipped.
func MountSSHFS(config config.VeshConfig) (ProcHandler, error) {
	log.Println("---[SSHFS Mount]---")

	debug.Println("Theoretical vault path:", config.VeraCryptVaultPath)

	if config.IsSSHFSMounted() {
		log.Println("SSHFS is already mounted")
		return nil, nil
	}
	handle := external.CreateSSHFS(config.ExtractSSHFSParams())
	err := handle.Start()
	if err != nil {
		return nil, fmt.Errorf("can't mount sshfs: %w", err)
	}

	return handle, nil
}

// MountVeraCrypt attempts to mount the VeraCrypt container specified in the config. If
// it is already mounted, it is skipped.
func MountVeraCrypt(config config.VeshConfig) (ProcHandler, error) {
	log.Println("---[VeraCrypt Mount]---")

	debug.Println("Theoretical repo path:", config.RepoPath)

	if config.IsVeraCryptMounted() {
		log.Println("VeraCrypt is already mounted")
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

	debug.Println("VeraCrypt path:", config.VeraCryptVaultPath)
	handle := external.CreateVeraCrypt(config.ExtractVeraCryptParams())
	err := handle.Start()
	if err != nil {
		return nil, err
	}

	return handle, nil
}

// DismountSSHFS dismounts SSHFS if the process handler is available.
func DismountSSHFS(handle ProcHandler) {
	if handle != nil {
		log.Println("---[SSHFS Dismount]---")
		handle.Stop()
	}
}

// DismountVeraCrypt dismounts VeraCrypt if the process handler is available.
func DismountVeraCrypt(handle ProcHandler) {
	if handle != nil {
		log.Println("---[VeraCrypt Dismount]---")
		handle.Stop()
	}
}

// Main implements the base logic of the application.
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

	utils.AddCleaning(func() error {
		DismountSSHFS(sshHandle)
		return nil
	}, "SSHFS")

	err = config.CheckSSHFS()
	if err != nil {
		return err
	}

	veraHandle, err := MountVeraCrypt(config)
	if err != nil {
		return err
	}
	utils.AddCleaning(func() error {
		DismountVeraCrypt(veraHandle)
		return nil
	}, "VeraCrypt")

	err = config.CheckVeraCrypt()
	if err != nil {
		return err
	}

	log.Println("---[Begin Main]---")
	err = helper.Loop(config)
	if err != nil {
		return err
	}
	log.Println("---[End Main]---")

	return nil
}

func main() {
	utils.InitiateCleaningState()
	signalHandler()
	defer utils.CleanStack()
	if err := Main(); err != nil {
		utils.CleanStack()
		log.Fatal(err)
	}
}
