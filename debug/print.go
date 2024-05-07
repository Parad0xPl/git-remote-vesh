package debug

import (
	"log"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

func Println(args ...any) {
	if utils.IsDebug() {
		log.Println(args...)
	}
}
func Printf(format string, args ...any) {
	if utils.IsDebug() {
		log.Printf(format, args...)
	}
}
