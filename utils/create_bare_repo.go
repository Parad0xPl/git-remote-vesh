package utils

import "os/exec"

// CreateBareRepo creates a bare Git repository using the "git init --bare" command.
func CreateBareRepo(path string) error {
	cmd := exec.Command("git", "init", "--bare", path)
	return cmd.Run()
}
