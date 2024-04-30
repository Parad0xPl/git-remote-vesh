package utils

import "os/exec"

func CreateBareRepo(path string) error {
	cmd := exec.Command("git", "init", "--bare", path)
	return cmd.Run()
}
