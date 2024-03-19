//go:build linux

package external

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

type SshfsWinHandle struct {
	startArguments []string
	mountPath      string
}

const defaultSSHFSPath = "/usr/bin/sshfs"

func getSSHFSPath() string {
	return getPathOrDef("sshfs", defaultSSHFSPath)
}

func CreateSSHFS(config config.SSHFSParams) *SshfsWinHandle {
	sshLogin := formatSSHConnection(&config)
	mountPath := config.SSHMountPath
	port := fmt.Sprintf("-p %d", config.SSHPort)

	ident_file := config.SSHIdentityFile

	arguments := []string{
		sshLogin,
		mountPath,
		port,
		"-o debug",
		"-o loglevel=debug1",
		"-o StrictHostKeyChecking=no",
		"-o UserKnownHostsFile=/dev/null",
		"-o large_read",
		"-o kernel_cache",
		"-o follow_symlinks",
	}

	if ident_file != "" {
		arguments = append(arguments,
			"-o PreferredAuthentications=publickey",
			fmt.Sprintf("-o IdentityFile=%s", ident_file),
		)
	}

	return &SshfsWinHandle{
		startArguments: arguments,
		mountPath:      mountPath,
	}
}

func (s *SshfsWinHandle) Start() error {
	path := getSSHFSPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find sshfs executable: %v", err)
	}
	cmd := exec.Command(path, s.startArguments...)

	if utils.IsDebug() {
		log.Println("---SSHFS Start---")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stderr
	}

	err := cmd.Run()
	log.Println("---SSHFS Start---")
	if err != nil {
		return fmt.Errorf("can't run sshfs executable: %v", err)
	}
	return nil
}

func (s *SshfsWinHandle) Stop() error {
	cmd := exec.Command("fusermount", "-u", s.mountPath)

	if utils.IsDebug() {
		log.Println("---SSHFS Stop---")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stderr
	}

	err := cmd.Run()
	log.Println("---SSHFS Stop---")
	if err != nil {
		return fmt.Errorf("can't unmount sshfs executable: %v", err)
	}
	return nil
}
