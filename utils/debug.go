package utils

import (
	"os"
	"strings"
)

// IsDebug checks if the DEBUG environment variable is set to "vesh". This enables
// additional logging and output for certain executed programs.
func IsDebug() bool {
	v := os.Getenv("DEBUG")
	v = strings.ToLower(v)
	if v == "vesh" {
		return true
	}
	return false
}
