package helper

import (
	"bytes"
	"fmt"
)

func (h *helperContext) getRefs() ([][2]string, error) {
	output, err := h.gitExec("for-each-ref",
		"--format=%(objectname) %(refname)",
		"refs/heads/")
	if err != nil {
		return nil, err
	}

	lines := bytes.Split(output, []byte{'\n'})

	num_of_lines := len(lines)
	out := make([][2]string, num_of_lines-1)
	i := 0

	for _, line := range lines {
		fields := bytes.SplitN(line, []byte{' '}, 2)
		if len(fields) < 2 {
			break
		}
		out[i] = [2]string{string(fields[0]), string(fields[1])}
		i++
	}

	return out, nil
}

func (h *helperContext) getHeadRef(name string) (string, error) {
	output, err := h.gitExec("symbolic-ref", name)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(output)), nil
}

func (h *helperContext) list() error {
	refs, err := h.getRefs()
	if err != nil {
		return err
	}

	head, err := h.getHeadRef("HEAD")
	if err != nil {
		return err
	}

	for _, ref := range refs {
		fmt.Printf("? %s\n", ref[1])
	}
	fmt.Printf("@%s HEAD\n", head)
	fmt.Println()
	return nil
}
