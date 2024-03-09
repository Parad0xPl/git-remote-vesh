package main

import (
	"log"
	"time"

	"github.com/Parad0xpl/git_enc/v2/windows"
)

type ProcHandler interface {
	Start() error
	Stop() error
}

func MountSSHFS() (ProcHandler, error) {
	log.Println("---[SSHFS Mount]---")

	handle := windows.CreateSSHFS()
	handle.Start()

	return handle, nil
}

func MountVeraCrypt() (ProcHandler, error) {
	log.Println("---[VeraCrypt Mount]---")

	handle := windows.CreateVeraCrypt()
	err := handle.Start()
	if err != nil {
		return nil, err
	}

	return handle, nil
}

func DismountSSHFS(handle ProcHandler) {
	log.Println("---[SSHFS Dismount]---")
	handle.Stop()
}

func DismountVeraCrypt(handle ProcHandler) {
	log.Println("---[VeraCrypt Dismount]---")
	handle.Stop()
}

func PushChanges() {
	time.Sleep(10 * time.Second)
}

func Main() error {
	GetConfig()

	sshHandle, err := MountSSHFS()
	if err != nil {
		return err
	}
	defer DismountSSHFS(sshHandle)
	veraHandle, err := MountVeraCrypt()
	if err != nil {
		return err
	}
	defer DismountVeraCrypt(veraHandle)

	PushChanges()

	return nil
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
