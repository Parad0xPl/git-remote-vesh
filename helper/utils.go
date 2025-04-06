package helper

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
)

// gitCmdPrepare prepares an exec.Cmd for running a Git command in the remote repository
// with the given arguments. It sets the GIT_DIR environment variable and returns a ready-to-run command.
func (h *helperContext) gitCmdPrepare(args ...string) *exec.Cmd {
	debug.Println("Git command:", args)
	cmd := exec.Command("git", args...)
	envs := cmd.Environ()
	envs = append(envs, fmt.Sprintf("GIT_DIR=%s", h.repoPath))
	cmd.Env = envs
	return cmd
}

// gitExec executes a Git command with the given arguments and returns the output as a byte array.
func (h *helperContext) gitExec(args ...string) ([]byte, error) {
	cmd := h.gitCmdPrepare(args...)
	return cmd.Output()
}

// gitExecStdout executes a Git command with output redirected to Stderr and Stdout.
func (h *helperContext) gitExecStdout(args ...string) error {
	cmd := h.gitCmdPrepare(args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// gitExecString executes a Git command with the given arguments and returns the output as a string.
func (h *helperContext) gitExecString(args ...string) (string, error) {
	output, err := h.gitExec(args...)
	return string(bytes.TrimSpace(output)), err
}

// ensureFile creates a file at the given path if it doesn't exist.
func ensureFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	debug.Println("Created file:", path)
	file.Close()
	return nil
}

// readCmd reads one line of input from Stdin and splits it into a slice of strings.
// It returns the split command or an error.
func (ctx *helperContext) readCmd() ([]string, error) {
	commandLine, err := ctx.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			debug.Printf("Got EOF - ignoring?")
			return nil, nil
		}
		return nil, fmt.Errorf("can't read the next line of communication: %v", err)
	}
	debug.Printf("Got command '%s'\n", commandLine)

	commandLineParts := strings.SplitN(commandLine, " ", 2)
	return commandLineParts, nil
}
