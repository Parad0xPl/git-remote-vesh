//go:build linux

package external

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
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

// CreateSSHFS creates an SSHFS connection handle for the current configuration.
// Returns a handle that can start and stop the resource.
func CreateSSHFS(config config.SSHFSParams) *SshfsWinHandle {
	sshLogin := formatSSHConnection(&config)
	mountPath := config.SSHMountPath
	port := fmt.Sprintf("-p %d", config.SSHPort)

	debug.Println("SSHFS remote address:", sshLogin)
	debug.Println("SSHFS mount path:", mountPath)

	ident_file := config.SSHIdentityFile

	os.MkdirAll(mountPath, 0o700)

	arguments := []string{
		sshLogin,
		mountPath,
		port,
		"-ologlevel=debug1",
		"-oStrictHostKeyChecking=no",
		"-oUserKnownHostsFile=/dev/null",
		"-okernel_cache",
		"-ofollow_symlinks",
		"-oidmap=user",
		"-oallow_root",
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

// Start initiates the SSHFS connection using the provided configuration.
func (s *SshfsWinHandle) Start() error {
	path := getSSHFSPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find sshfs executable: %v", err)
	}
	cmd := exec.Command(path, s.startArguments...)

	if utils.IsDebug() {
		log.Println("---SSHFS Start---")
		log.Println("Arguments:", s.startArguments)
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

// Stop terminates the SSHFS connection and unmounts the resource.
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
