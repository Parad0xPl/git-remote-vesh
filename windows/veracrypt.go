package windows

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

type VeraCryptWinHandle struct {
	mountLetter string
	vaultPath   string
}

const defautVeraCryptPath = "C:\\Program Files\\VeraCrypt\\VeraCrypt.exe"

func getVeraCryptExecPath() string {
	path, _ := exec.LookPath("VeraCrypt")

	if path == "" {
		path = defautVeraCryptPath
	}

	return path
}

func CreateVeraCrypt(config config.VeraCryptParams) *VeraCryptWinHandle {
	return &VeraCryptWinHandle{
		mountLetter: config.VeraCryptMountPath,
		vaultPath:   config.VeraCryptVaultPath,
	}
}

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
