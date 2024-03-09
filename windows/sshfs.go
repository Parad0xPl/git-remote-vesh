package windows

import (
	"fmt"
	"os/exec"
)

type SshfsWinHandle struct {
	args []string
	cmd  *Cmd
}

const defaultPath = "C:\\Program Files\\SSHFS-Win\\bin\\sshfs.exe"

func getExecPath() string {
	path, _ := exec.LookPath("sshfs")

	if path == "" {
		path = defaultPath
	}

	return path
}

func createWinSSHFS() *SshfsWinHandle {
	SSHName := "testgit"
	SSHHost := "parad0x.pl"
	SSHPath := "/home/but_revolver"
	ssh_login := fmt.Sprintf("%s@%s:%s", SSHName, SSHHost, SSHPath)
	mount_letter := "Z:"
	port := fmt.Sprintf("-p%d", 22)

	ident_file = "~/.ident"

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
			fmt.Sprintf("-oIdentitiyFile=\"%s\"", ident_file))
	}

	return &SshfsWinHandle{
		args: arguments,
	}
}

func (s *SshfsWinHandle) Start() error {
	path := getExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find sshfs executable: %s", err)
	}
	s.cmd = exec.Command(path, s.args...)
	err = s.Start()
	if err != nil {
		return fmt.Errorf("can't start sshfs executable: %s", err)
	}
	return nil
}

func (s *SshfsWinHandle) Stop() error {
	return nil
}
