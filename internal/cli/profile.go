package cli

import (
	"flag"
	"fmt"

	"github.com/hamajun504/ProfilesCLI/internal/profile"
)

func runProfile(args []string) error {
	fs := flag.NewFlagSet("profile create", flag.ContinueOnError)

	name := fs.String("name", "", "profile name")
	user := fs.String("user", "", "user name")
	project := fs.String("project", "", "project name")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	switch args[0] {
	case "create":
		fmt.Println("Create Profile")
		profile.Create(*name, *user, *project)
	case "get":
		fmt.Println("Return Profile")
	case "list":
		fmt.Println("List all profiles")
	case "delete":
		fmt.Println("Delete Profile")
	}
	return nil
}
