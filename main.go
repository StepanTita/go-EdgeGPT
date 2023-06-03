package main

import (
	"os"

	"github.com/StepanTita/go-EdgeGPT/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(2)
	}
	os.Exit(0)
}
