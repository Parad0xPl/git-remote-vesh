//go:build windows

package external

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

// VeraCryptWinHandle represents a handle for managing VeraCrypt operations on Windows.
type VeraCryptWinHandle struct {
	mountLetter string
	vaultPath   string
}

const defaultVeraCryptPath = "C:\\Program Files\\VeraCrypt\\VeraCrypt.exe"

// getVeraCryptExecPath retrieves the path to the VeraCrypt executable,
// or returns the default path if not found in the system PATH.
func getVeraCryptExecPath() string {
	return getPathOrDef("VeraCrypt", defaultVeraCryptPath)
}

// CreateVeraCrypt creates a VeraCrypt handle using the provided configuration.
func CreateVeraCrypt(config config.VeraCryptParams) *VeraCryptWinHandle {
	return &VeraCryptWinHandle{
		mountLetter: config.VeraCryptMountPath,
		vaultPath:   config.VeraCryptVaultPath,
	}
}

// Start mounts the VeraCrypt volume using the provided configuration.
func (s *VeraCryptWinHandle) Start() error {
	path := getVeraCryptExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %v", err)
	}

	arguments :=
		[]string{
			"/q",
			"/l",
			s.mountLetter,
			"/a",
			"/v",
			s.vaultPath,
		}

	cmd := exec.Command(path, arguments...)
	_, err := cmd.Output()
	// log.Printf("Output: %s\n", string(output))
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %v", err)
	}
	return nil
}

// Stop unmounts the VeraCrypt volume.
func (s *VeraCryptWinHandle) Stop() error {
	path := getVeraCryptExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %v", err)
	}

	arguments :=
		[]string{
			"/q",
			"/d",
			s.mountLetter,
		}

	cmd := exec.Command(path, arguments...)
	_, err := cmd.Output()
	// log.Println(string(output))
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %v", err)
	}
	return nil
}
