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
		profile.Create(*name, *user, *project)

	case "get":
		user, project, err := profile.Get(*name)
		if err != nil {
			return err
		}
		fmt.Println(user, project)

	case "list":
		profiles, err := profile.List()
		if err != nil {
			return err
		}
		printProfiles(profiles)

	case "delete":
		if err := profile.Delete(*name); err != nil {
			return err
		}
	}
	return nil
}

func printProfiles(profiles []string) {
	for _, name := range profiles {
		fmt.Println(name)
	}
}
