package windows

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

type VeraCryptHandle struct {
	mountPath string
	vaultPath string
}

const defautVeraCryptPath = "/usr/bin/veracrypt"

func getVeraCryptExecPath() string {
	path, _ := exec.LookPath("veracrypt")

	if path == "" {
		path = defautVeraCryptPath
	}

	return path
}

func CreateVeraCrypt(config config.VeraCryptParams) *VeraCryptHandle {
	return &VeraCryptHandle{
		mountPath: config.VeraCryptMountPath,
		vaultPath: config.VeraCryptVaultPath,
	}
}

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

	cmd := exec.Command(path, arguments...)
	err := cmd.Run()
	// log.Printf("Output: %s\n", string(output))
	if err != nil {
		return fmt.Errorf("can't start veracrypt executable: %v", err)
	}
	return nil
}

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
