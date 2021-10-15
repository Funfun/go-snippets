package main

import (
	"os"

	"github.com/Funfun/go-snippets/go-elk/go-pubsub-elk/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
