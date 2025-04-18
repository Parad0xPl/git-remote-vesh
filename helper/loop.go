package helper

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Parad0xpl/git-remote-vesh/v2/config"
)

// capabilities prints the capabilities handled by our helper.
func (h *helperContext) capabilities() {
	fmt.Println("import")
	fmt.Println("export")
	fmt.Printf("refspec %s\n", h.headRefspec)
	fmt.Printf("*import-marks %s\n", h.gitmarks)
	fmt.Printf("*export-marks %s\n", h.gitmarks)
	fmt.Println()
}

// helperContext gathers data required between all function calls.
type helperContext struct {
	headRefspec string
	repoPath    string
	gitmarks    string
	veshmarks   string

	reader *bufio.Reader
}

// prepare creates all necessary files inside the local repository.
func prepare(config config.VeshConfig) (helperContext, error) {
	veshDir := config.GetVeshConfigDir()
	ctx := helperContext{
		headRefspec: fmt.Sprintf("refs/heads/*:refs/remotes/%s/heads/*", config.RemoteName),
		repoPath:    config.RepoPath,
		gitmarks:    filepath.Join(veshDir, "git.marks"),
		veshmarks:   filepath.Join(veshDir, "veshgit.marks"),
	}

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

// Loop contains the full logic for handling Git communication. Most of the logic
// is based on Alec Newman's article (https://rovaughn.github.io/2015-2-9.html).
func Loop(config config.VeshConfig) error {
	ctx, err := prepare(config)
	if err != nil {
		return fmt.Errorf("can't prepare vesh context: %v", err)
	}

	ctx.reader = bufio.NewReader(os.Stdin)
	for {
		commandLineParts, err := ctx.readCmd()
		if err != nil {
			return err
		}
		if commandLineParts == nil {
			return nil
		}

		command := commandLineParts[0]

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
			return fmt.Errorf("unknown command '%s' from Git", command)
		}
	}
}
