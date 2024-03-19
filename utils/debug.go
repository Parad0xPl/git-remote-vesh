package utils

import (
	"os"
	"strings"
)

func IsDebug() bool {
	v := os.Getenv("DEBUG")
	v = strings.ToLower(v)
	if v == "vesh" {
		return true
	}
	return false
}
