package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/hamajun504/ProfilesCLI/internal/profile"
)

func runProfile(args []string) error {
	fs := flag.NewFlagSet("profile create", flag.ContinueOnError)

	name := fs.String("name", "", "profile name")
	user := fs.String("user", "", "user name")
	project := fs.String("project", "", "project name")

	if len(args) == 0 || len(args) == 1 {
		printHelp()
		return nil
	}

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	switch args[0] {
	case "create":
		exist, err := profile.Exist(*name)
		if err != nil {
			return err
		}
		if exist && !AskToOverwrite(*name) {
			return nil
		}
		if err := profile.Create(*name, *user, *project); err != nil {
			return err
		}

	case "get":
		user, project, err := profile.Get(*name)
		if err != nil {
			return err
		}
		printGetOutput(*name, user, project)

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

	case "help":
		printHelp()
	}
	return nil
}

func printProfiles(profiles []string) {
	for _, name := range profiles {
		fmt.Println(name)
	}
}

func printHelp() error {
	_, err := fmt.Println(help_message)
	return err
}

func printGetOutput(name, user, project string) error {
	output := "profile:  " + name + "\n" +
		"user   :  " + user + "\n" +
		"project:  " + project
	_, err := fmt.Println(output)
	return err
}

func AskToOverwrite(name string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Profile %s already exist. Do you really want to overwrite it? [y/[n]]", name)
	answer, _ := reader.ReadString('\n')
	// err != nil means input ended not with \n, thats acceptable

	answer = strings.TrimSpace(strings.ToLower(answer))

	return answer == "y" || answer == "yes"
}
