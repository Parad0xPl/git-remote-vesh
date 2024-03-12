package helper

import (
	"fmt"
	"os"
)

func (h *helperContext) export() error {
	beforeRefs, err := h.getRefs()
	if err != nil {
		return fmt.Errorf("command export: can't get before ref list: %v", err)
	}

	cmd := h.gitCmdPrepare("fast-import", "--quiet",
		"--import-marks="+h.veshmarks,
		"--export-marks="+h.veshmarks,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command export: git fast-import: %v", err)
	}

	afterRefs, err := h.getRefs()
	beforeRefsMap := transformSlice(beforeRefs)
	if err != nil {
		return fmt.Errorf("command export: can't get after ref list: %v", err)
	}

	for _, el := range afterRefs {
		refname := el[1]
		if beforeRefsMap[refname] != el[0] {
			fmt.Printf("ok %s\n", refname)
		}
	}

	fmt.Printf("\n")

	return nil
}
