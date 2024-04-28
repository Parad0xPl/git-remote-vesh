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
	if utils.IsDebug() {
		log.Println("Created file:", path)
	}
	file.Close()
	return nil
}

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
