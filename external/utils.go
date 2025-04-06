package external

import (
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

// getPathOrDef tries to find the path to the executable and
// returns the default path if not found in the system PATH.
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

// formatSSHConnection prepares the SSH connection address from the configuration.
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
