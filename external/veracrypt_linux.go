//go:build linux

package external

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

// VeraCryptHandle represents a handle for managing VeraCrypt operations.
type VeraCryptHandle struct {
	mountPath string
	vaultPath string
}

const defautVeraCryptPath = "/usr/bin/veracrypt"

// getVeraCryptExecPath retrieves the path to the VeraCrypt executable,
// or returns the default path if not found in the system PATH.
func getVeraCryptExecPath() string {
	return getPathOrDef("veracrypt", defautVeraCryptPath)
}

// CreateVeraCrypt creates a VeraCrypt handle using the provided configuration.
func CreateVeraCrypt(config config.VeraCryptParams) *VeraCryptHandle {
	return &VeraCryptHandle{
		mountPath: config.VeraCryptMountPath,
		vaultPath: config.VeraCryptVaultPath,
	}
}

// Start mounts the VeraCrypt volume using the provided configuration.
func (s *VeraCryptHandle) Start() error {
	path := getVeraCryptExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %v", err)
	}

	arguments :=
		[]string{
			s.vaultPath,
			s.mountPath,
		}

	os.MkdirAll(s.mountPath, 0o700)

	cmd := exec.Command(path, arguments...)
	err := cmd.Run()
	// log.Printf("Output: %s\n", string(output))
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %v", err)
	}
	return nil
}

// Stop unmounts the VeraCrypt volume.
func (s *VeraCryptHandle) Stop() error {
	path := getVeraCryptExecPath()

	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("can't find veracrypt executable: %v", err)
	}

	arguments :=
		[]string{
			"-d",
			s.mountPath,
		}

	cmd := exec.Command(path, arguments...)
	err := cmd.Run()
	// log.Println(string(output))
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %v", err)
	}
	return nil
}
