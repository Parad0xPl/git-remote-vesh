package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

func singalHandler() {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		log.Println("Got signal:", sig)
		utils.CleanStack()
		os.Exit(0)
	}()

}
