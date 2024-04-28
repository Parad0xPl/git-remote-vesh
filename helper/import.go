package helper

import (
	"fmt"
	"log"
	"strings"
)

func (h *helperContext) import_(f string) error {
	refs := make([]string, 0)
	for {
		log.Printf("Importing file: %s", f)

		ref := strings.TrimSpace(f)
		refs = append(refs, ref)

		commandLineParts, err := h.readCmd()
		if err != nil {
			return err
		}
		if commandLineParts[0] != "import" {
			break
		}
		f = commandLineParts[1]
	}

	fmt.Printf("feature import-marks=%s\n", h.gitmarks)
	fmt.Printf("feature export-marks=%s\n", h.gitmarks)
	fmt.Printf("feature done\n")

	args := []string{"fast-export",
		"--import-marks=" + h.veshmarks,
		"--export-marks=" + h.veshmarks,
		"--refspec=" + h.headRefspec,
	}

	args = append(args, refs...)

	log.Println(args)
	err := h.gitExecStdout(args...)
	if err != nil {
		return fmt.Errorf("import fail: %v", err)
	}

	fmt.Printf("done\n")

	return nil
}
