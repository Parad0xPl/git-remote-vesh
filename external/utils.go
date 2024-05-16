package external

import (
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

// getPathOrDef tries to find path to the executable and
// return default if not find in PATH
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

// formatSSHConnection prepare ssh connection address from config
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
	} else if runtime.GOOS == "linux" {
		output = output + ":"
	}

	return output
}
