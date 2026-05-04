package main

import (
	"os"

	"github.com/hamajun504/ProfilesCLI/internal/cli"
)

func main() {
	cli.Run(os.Args[1:])
}
