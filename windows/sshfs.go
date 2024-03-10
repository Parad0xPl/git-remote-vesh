package windows

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Parad0xpl/git_enc/v2/config"
)

type SshfsWinHandle struct {
	args []string
	cmd  *exec.Cmd
}

const defaultPath = "C:\\Program Files\\SSHFS-Win\\bin\\sshfs.exe"

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
	port := fmt.Sprintf("-p%d", 22)

	ident_file := config.SSHIdentityFile

	arguments := []string{
		ssh_login,
		mount_letter,
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
		return fmt.Errorf("can't find sshfs executable: %s", err)
	}
	s.cmd = exec.Command(path, s.args...)
	// s.cmd.Stdout = os.Stdout
	// s.cmd.Stderr = os.Stderr

	s.cmd.Env = []string{
		fmt.Sprintf("PATH=%s", filepath.Dir(path)),
	}
	err := s.cmd.Start()
	if err != nil {
		return fmt.Errorf("can't start sshfs executable: %s", err)
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
