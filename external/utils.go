package external

import (
	"os/exec"
	"path/filepath"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

func getPathOrDef(name, defPath string) string {
	path, _ := exec.LookPath(name)

	if path == "" {
		path = defPath
	}

	path2, err := filepath.Abs(path)

	if err == nil {
		path = path2
	}

	return path
}

func formatSSHConnection(config *config.SSHFSParams) string {
	SSHName := config.SSHUser
	SSHHost := config.SSHAddress
	SSHPath := config.SSHRemotePath

	output := SSHHost

	if SSHName != "" {
		output = SSHName + "@" + output
	}

	if SSHPath != "" {
		output = output + ":" + SSHPath
	}

	return output
}
