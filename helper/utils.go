package helper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// gitCmdPrepare prepare exec.Cmd for running git command in remote repository
// with given arguments. It sets GIT_DIR envirenment variable and return ready
// to run.
func (h *helperContext) gitCmdPrepare(args ...string) *exec.Cmd {
	if utils.IsDebug() {
		log.Println("Git command:", args)
	}
	cmd := exec.Command("git", args...)
	envs := cmd.Environ()
	envs = append(envs, fmt.Sprintf("GIT_DIR=%s", h.repoPath))
	cmd.Env = envs
	return cmd
}

// gitExec execute git with given arguments. Return output as bytes array
func (h *helperContext) gitExec(args ...string) ([]byte, error) {
	cmd := h.gitCmdPrepare(args...)
	return cmd.Output()
}

// gitExecStdout execute git command with output redirected to Stderr and Stdout
func (h *helperContext) gitExecStdout(args ...string) error {
	cmd := h.gitCmdPrepare(args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// gitExecString execute git with given argument. Return output as string
func (h *helperContext) gitExecString(args ...string) (string, error) {
	output, err := h.gitExec(args...)
	return string(bytes.TrimSpace(output)), err
}

// ensureFile create file with given path it doesn't exists
func ensureFile(path string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	if utils.IsDebug() {
		log.Println("Created file:", path)
	}
	file.Close()
	return nil
}

// readCmd return one line of command taken from Stdin. Return splitted by space
// slice of strings
func (ctx *helperContext) readCmd() ([]string, error) {
	commandLine, err := ctx.reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			if utils.IsDebug() {
				log.Printf("Got EOF - ignoring?")
			}
			return nil, nil
		}
		return nil, fmt.Errorf("can't read next line of communication: %v", err)
	}
	if utils.IsDebug() {
		log.Printf("Got command '%s'\n", commandLine)
	}

	commandLineParts := strings.SplitN(commandLine, " ", 2)
	return commandLineParts, nil
}
