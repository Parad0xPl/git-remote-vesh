package main

import (
	"log"
	"os"

	"github.com/Parad0xpl/git_enc/v2/config"
	"github.com/Parad0xpl/git_enc/v2/helper"
	"github.com/Parad0xpl/git_enc/v2/windows"
)

type ProcHandler interface {
	Start() error
	Stop() error
}

func MountSSHFS(config config.EncConfig) (ProcHandler, error) {
	log.Println("---[SSHFS Mount]---")

	if _, err := os.Stat(config.VeraCryptVaultPath); err == nil {
		log.Println("Vault already found - skipping sshfs mount")
		return nil, nil
	}
	handle := windows.CreateSSHFS(config.ExtractSSHFSParams())
	handle.Start()

	return handle, nil
}

func MountVeraCrypt(config config.EncConfig) (ProcHandler, error) {
	log.Println("---[VeraCrypt Mount]---")

	if _, err := os.Stat(config.RepoPath); err == nil {
		log.Println("Repo already found - skipping VeraCrypt mount")
		return nil, nil
	}

	handle := windows.CreateVeraCrypt(config.ExtractVeraCryptParams())
	err := handle.Start()
	if err != nil {
		return nil, err
	}

	return handle, nil
}

func DismountSSHFS(handle ProcHandler) {
	log.Println("---[SSHFS Dismount]---")
	if handle != nil {
		handle.Stop()
	}
}

func DismountVeraCrypt(handle ProcHandler) {
	log.Println("---[VeraCrypt Dismount]---")
	if handle != nil {
		handle.Stop()
	}
}

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

	err = helper.Loop(config)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
