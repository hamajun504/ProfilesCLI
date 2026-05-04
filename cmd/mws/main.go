package main

import (
	"os"

	"github.com/hamajun504/cmd/mws/internal/cli"
)

func main() {
	cli.Run(os.Args[1:])
}
