package main

import (
	"os"

	"github.com/things-go/ens/cmd/ormat/command"
)

func main() {
	err := command.NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
