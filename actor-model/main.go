package main

import (
	"os"
	"os/signal"

	"github.com/Funfun/go-snippets/actor-model/actor"
)

func main() {
	actor := actor.Actor{}
	go actor.Process()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	actor.Stop()
}
