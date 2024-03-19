package helper

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

func (h *helperContext) capabilities() {
	// log.Println("--Printing cap--")
	fmt.Println("import")
	fmt.Println("export")
	fmt.Printf("refspec %s\n", h.headRefspec)
	fmt.Printf("refspec %s\n", h.tagRefspec)
	fmt.Printf("*import-marks %s\n", h.gitmarks)
	fmt.Printf("*export-marks %s\n", h.gitmarks)
	fmt.Println()
	// log.Println("--Printed cap--")
}

type helperContext struct {
	headRefspec string
	tagRefspec  string
	repoPath    string
	gitmarks    string
	veshmarks   string

	reader *bufio.Reader
}

func prepare(config config.EncConfig) (helperContext, error) {
	veshDir := config.GetVeshConfigDir()
	ctx := helperContext{
		headRefspec: fmt.Sprintf("refs/heads/*:refs/vesh/%s/heads/*", config.RemoteName),
		tagRefspec:  fmt.Sprintf("refs/tags/*:refs/vesh/%s/tags/*", config.RemoteName),
		repoPath:    config.RepoPath,
		gitmarks:    filepath.Join(veshDir, "git.marks"),
		veshmarks:   filepath.Join(veshDir, "veshgit.marks"),
	}

	// log.Printf("Repo path: %s\n", ctx.repoPath)
	err := os.MkdirAll(veshDir, 0755)
	if err != nil {
		return helperContext{}, err
	}

	err = ensureFile(ctx.gitmarks)
	if err != nil {
		return helperContext{}, err
	}
	err = ensureFile(ctx.veshmarks)
	if err != nil {
		return helperContext{}, err
	}

	return ctx, nil
}

// / Loop contain full logic of handling git communication. Most of logic is
// / based on Alec Newman article (https://rovaughn.github.io/2015-2-9.html)
func Loop(config config.EncConfig) error {
	ctx, err := prepare(config)
	if err != nil {
		return fmt.Errorf("can't prepare vesh contex: %v", err)
	}

	// log.Println("---[Main helper loop]---")
	ctx.reader = bufio.NewReader(os.Stdin)
	for {
		// log.Println("--- loop run ---")
		commandLineParts, err := ctx.readCmd()
		if err != nil {
			return err
		}

		command := commandLineParts[0]
		// log.Printf("command first part '%s'\n", command)

		if command == "\n" {
			return nil
		} else if command == "capabilities\n" {
			ctx.capabilities()
		} else if command == "list\n" {
			ctx.list()
		} else if command == "export\n" {
			ctx.export()
		} else if command == "import" {
			ctx.import_(commandLineParts[1])
		} else {
			return fmt.Errorf("unkown command '%s' from git", command)
		}
	}
}
