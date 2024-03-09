package main

import (
	"log"
)

type ProcHandler interface {
	Start() error
	Stop() error
}

func MountSSHFS() ProcHandler {
	log.Println("---[SSHFS Mount]---")

	return nil
}

func MountVeraCrypt() ProcHandler {
	log.Println("---[VeraCrypt Mount]---")

	return nil
}

func DismountSSHFS(handle ProcHandler) {
	log.Println("---[SSHFS Dismount]---")
	handle.Stop()
}

func DismountVeraCrypt(handle ProcHandler) {
	log.Println("---[VeraCrypt Dismount]---")
	handle.Stop()
}

func PushChanges() {}

func Main() error {
	GetConfig()

	sshHandle := MountSSHFS()
	veraHandle := MountVeraCrypt()

	PushChanges()

	DismountVeraCrypt(veraHandle)
	DismountSSHFS(sshHandle)

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
