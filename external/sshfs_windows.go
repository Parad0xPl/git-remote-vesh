//go:build windows

package external

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

type SshfsWinHandle struct {
	args []string
	cmd  *exec.Cmd
}

const defaultSSHFSPath = "C:\\Program Files\\SSHFS-Win\\bin\\sshfs.exe"

func getSSHFSPath() string {
	return getPathOrDef("sshfs", defaultSSHFSPath)
}

func CreateSSHFS(config config.SSHFSParams) *SshfsWinHandle {
	sshLogin := formatSSHConnection(&config)
	if config.SSHRemotePath == "" {
		sshLogin = sshLogin + ":"
	}
	mountLetter := config.SSHMountPath
	port := fmt.Sprintf("-p%d", config.SSHPort)

	identFile := config.SSHIdentityFile

	arguments := []string{
		sshLogin,
		mountLetter,
		port,
		"-ovolname=BitEncryptSSHFS",
		"-odebug",
		"-ologlevel=debug1",
		"-oStrictHostKeyChecking=no",
		"-oUserKnownHostsFile=/dev/null",
		"-oidmap=user",
		"-ouid=-1",
		"-ogid=-1",
		"-oumask=000",
		"-ocreate_umask=000",
		"-omax_readahead=1GB",
		"-oallow_other",
		"-olarge_read",
		"-okernel_cache",
		"-ofollow_symlinks",
	}

	debug.Println("Arguments for sshfs:", arguments)

	if identFile != "" {
		arguments = append(arguments,
			"-oPreferredAuthentications=publickey",
			fmt.Sprintf("-oIdentityFile=%s", identFile),
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

	s.cmd.Env = []string{
		fmt.Sprintf("PATH=%s", filepath.Dir(path)),
	}

	if utils.IsDebug() {
		log.Println("SSHFS Env:", s.cmd.Env)

		log.Println("---SSHS Start---")
		s.cmd.Stderr = os.Stderr
		s.cmd.Stdout = os.Stderr
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
		debug.Println("Killing SSHFS process")
		s.cmd.Process.Kill()
	} else {
		debug.Println("SSHFS Process is not available")
	}
	return nil
}
