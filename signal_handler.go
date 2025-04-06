package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
)

// signalHandler sets up a signal handler to clean up resources and exit gracefully
// when a termination signal is received.
// It listens for SIGINT, SIGHUP, SIGTERM, and SIGQUIT signals.
// When a signal is received, it logs the signal and calls CleanStack to clean up resources.
// Finally, it exits the program with a status code of 0.
// This function runs in a separate goroutine to avoid blocking the main program.
func signalHandler() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sigs
		log.Println("Got signal:", sig)
		utils.CleanStack()
		os.Exit(0)
	}()

}
