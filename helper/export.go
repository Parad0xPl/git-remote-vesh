package helper

import (
	"fmt"
	"os"
)

// export handles all logic of exporting files from the local repository to the remote.
func (h *helperContext) export() error {
	// Save the list of refs before running `git fast-import` to compare
	// later what has changed.
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

	// Get the list of refs after the import.
	afterRefs, err := h.getRefs()
	beforeRefsMap := transformSlice(beforeRefs)
	if err != nil {
		return fmt.Errorf("command export: can't get after ref list: %v", err)
	}

	// Report "ok" for each changed ref.
	for _, el := range afterRefs {
		// log.Printf("Looking for match element %s = %s", el[1], el[0])
		refname := el[1]
		if beforeRefsMap[refname] != el[0] {
			// log.Printf(" - Found!!! [%s]", beforeRefsMap[refname])
			fmt.Printf("ok %s\n", refname)
		}
	}

	fmt.Printf("\n")

	return nil
}
