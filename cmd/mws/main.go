package main

import (
	"os"

	"github.com/hamajun504/mws/internal/cli"
)

func main() {
	cli.Run(os.Args[1:])
}
