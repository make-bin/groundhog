package main

import (
	"os"

	"github.com/make-bin/groundhog/pkg/interface/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}
