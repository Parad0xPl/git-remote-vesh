package utils

import (
	"os"
	"strings"
)

// IsDebug checks if debug variable is set to vesh. This enables mostly
// additional logging and output of few executed programs.
func IsDebug() bool {
	v := os.Getenv("DEBUG")
	v = strings.ToLower(v)
	if v == "vesh" {
		return true
	}
	return false
}
