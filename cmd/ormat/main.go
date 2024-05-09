package main

import (
	"os"

	"github.com/things-go/ens/cmd/ormat/command"
)

var root = command.NewRootCmd()

func main() {
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
