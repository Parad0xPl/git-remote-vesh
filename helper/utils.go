package helper

import (
	"bytes"
	"fmt"
	"os/exec"
)

func (h *helperContext) gitCmdPrepare(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	envs := cmd.Environ()
	envs = append(envs, fmt.Sprintf("GIT_DIT=%s", h.repoPath))
	cmd.Env = envs
	return cmd
}

func (h *helperContext) gitExec(args ...string) ([]byte, error) {
	cmd := h.gitCmdPrepare(args...)
	return cmd.Output()
}

func (h *helperContext) gitExecString(args ...string) (string, error) {
	output, err := h.gitExec(args...)
	return string(bytes.TrimSpace(output)), err
}
