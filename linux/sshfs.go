package windows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

type SshfsWinHandle struct {
	args []string
	cmd  *exec.Cmd
}

const defaultPath = "/usr/bin/sshfs"

func getSSHFSPath() string {
	path, _ := exec.LookPath("sshfs")

	if path == "" {
		path = defaultPath
	}

	path2, err := filepath.Abs(path)

	if err == nil {
		path = path2
	}

	return path
}

func CreateSSHFS(config config.SSHFSParams) *SshfsWinHandle {
	SSHName := config.SSHUser
	SSHHost := config.SSHAddress
	SSHPath := config.SSHRemotePath
	ssh_login := fmt.Sprintf("%s@%s:%s", SSHName, SSHHost, SSHPath)
	mount_letter := config.SSHMountPath
	port := fmt.Sprintf("-p %d", config.SSHPort)

	ident_file := config.SSHIdentityFile

	arguments := []string{
		ssh_login,
		mount_letter,
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
			"-oPreferredAuthentications=publickey",
			fmt.Sprintf("-oIdentityFile=%s", ident_file),
		)
	}

	return &SshfsWinHandle{
		args: arguments,
	}
}

func (s *SshfsWinHandle) Start() error {
	path := getSSHFSPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find sshfs executable: %v", err)
	}
	s.cmd = exec.Command(path, s.args...)
	// s.cmd.Stdout = os.Stdout
	// s.cmd.Stderr = os.Stderr

	s.cmd.Env = []string{
		fmt.Sprintf("PATH=%s", filepath.Dir(path)),
	}
	err := s.cmd.Start()
	if err != nil {
		return fmt.Errorf("can't start sshfs executable: %v", err)
	}
	time.Sleep(time.Second * 2)
	//TODO Check for error
	return nil
}

func (s *SshfsWinHandle) Stop() error {
	if s.cmd.Process != nil {
		s.cmd.Process.Kill()
	}
	return nil
}
