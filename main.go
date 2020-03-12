package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/corpulent/nuxx/cmd"
)

func init() {
}

func cleanup() {
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	cmd.Execute()
}
