package helper

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func (h *helperContext) gitCmdPrepare(args ...string) *exec.Cmd {
	cmd := exec.Command("git", args...)
	envs := cmd.Environ()
	log.Printf("Git repo dir: %s\n", h.repoPath)
	envs = append(envs, fmt.Sprintf("GIT_DIR=%s", h.repoPath))
	cmd.Env = envs
	return cmd
}

func (h *helperContext) gitExec(args ...string) ([]byte, error) {
	cmd := h.gitCmdPrepare(args...)
	return cmd.Output()
}

func (h *helperContext) gitExecStdout(args ...string) error {
	cmd := h.gitCmdPrepare(args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (h *helperContext) gitExecString(args ...string) (string, error) {
	output, err := h.gitExec(args...)
	return string(bytes.TrimSpace(output)), err
}

func ensureFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	file.Close()
	return nil
}

func (ctx *helperContext) readCmd() ([]string, error) {
	commandLine, err := ctx.reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("can't read next line of communication: %v", err)
	}
	// log.Printf("Got command '%s'\n", command_line)

	commandLineParts := strings.SplitN(commandLine, " ", 2)
	return commandLineParts, nil
}
