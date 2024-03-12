package helper

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Parad0xpl/git_enc/v2/config"
)

func (h *helperContext) capabilities() {
	fmt.Println("import")
	fmt.Println("export")
	fmt.Printf("refspec %s\n", h.headRefspec)
	fmt.Printf("refspec %s\n", h.tagRefspec)
	fmt.Printf("*import-marks %s", h.gitmarks)
	fmt.Printf("*export-marks %s", h.gitmarks)
	fmt.Println()
}

func (h *helperContext) export() {

}

func (h *helperContext) import_() {

}

type helperContext struct {
	headRefspec string
	tagRefspec  string
	repoPath    string
	gitmarks    string
	veshmarks   string
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

func prepare(config config.EncConfig) (helperContext, error) {
	ctx := helperContext{
		headRefspec: fmt.Sprintf("refs/heads/*:refs/vesh/%s/heads/*", config.RemoteName),
		tagRefspec:  fmt.Sprintf("refs/tags/*:refs/vesh/%s/tags/*", config.RemoteName),
		repoPath:    config.RepoPath,
		gitmarks:    path.Join(config.RepoPath, "git.marks"),
		veshmarks:   path.Join(config.RepoPath, "veshgit.marks"),
	}
	ensureFile(ctx.gitmarks)
	ensureFile(ctx.veshmarks)

	return ctx, nil
}

// / Loop contain full logic of handling git communication. Most of logic is
// / based on Alec Newman article (https://rovaughn.github.io/2015-2-9.html)
func Loop(config config.EncConfig) error {
	ctx, err := prepare(config)
	if err != nil {
		return fmt.Errorf("can't prepare vesh contex: %v", err)
	}

	log.Println("---[Main helper loop]---")
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		command_line, err := stdinReader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("can't read next line of communication: %v", err)
		}
		log.Printf("Got command '%s'\n", command_line)

		command_line_parts := strings.Split(command_line, " ")

		command := command_line_parts[0]

		if command == "\n" {
			return nil
		} else if command == "capabilities\n" {
			ctx.capabilities()
		} else if command == "list\n" {
			ctx.list()
		} else if command == "export\n" {
			ctx.export()
		} else if command == "import\n" {
			ctx.import_()
		} else {
			return fmt.Errorf("unkown command '%s' from git", command)
		}
	}
}
