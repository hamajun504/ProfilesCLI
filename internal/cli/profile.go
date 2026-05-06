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
	extendedFiles := fs.Bool("e", false, "output files with extra fields")
	allFiles := fs.Bool("a", false, "output all yaml-files")
	longOutput := fs.Bool("l", false, "detailed output")
	forceOverwrite := false
	fs.BoolVar(&forceOverwrite, "force", false, "overwrite existing profile")
	fs.BoolVar(&forceOverwrite, "f", false, "overwrite existing profile")

	if len(args) == 0 {
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
		if !exist {
			if err := profile.Create(*name, *user, *project); err != nil {
				return err
			}
			return nil
		}
		if !forceOverwrite {
			ow, err := askToOverwrite(*name)
			if err != nil {
				return err
			}
			if !ow {
				return nil
			}
		}
		if err := profile.Update(*name, *user, *project); err != nil {
			return err
		}

	case "get":
		p, err := profile.Get(*name)
		if err != nil {
			return err
		}
		if err := printGetOutput(p); err != nil {
			return err
		}

	case "list":
		var err error
		var profiles []profile.Profile
		if *allFiles {
			profiles, err = profile.List(profile.All)
		} else if *extendedFiles {
			profiles, err = profile.List(profile.ValidOrExtended)
		} else {
			profiles, err = profile.List(profile.Valid)
		}
		if err != nil {
			return err
		}
		if *longOutput {
			if err := printProfilesDetails(profiles); err != nil {
				return err
			}
		} else {
			if err := printProfilesShortly(profiles); err != nil {
				return err
			}
		}

	case "delete":
		exist, err := profile.Exist(*name)
		if err != nil {
			return err
		}
		if !exist {
			fmt.Fprintln(os.Stderr, "profile not exist")
		}
		if err := profile.Delete(*name); err != nil {
			return err
		}

	case "help":
		if err := printHelp(); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown profile command: %s", args[0])
	}

	return nil
}

func printProfilesShortly(profiles []profile.Profile) error {
	for _, p := range profiles {
		_, err := fmt.Println(p.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func printProfilesDetails(profiles []profile.Profile) error {
	var namesMaxLen, usersMaxLen, projectsMaxLen int
	for i := range profiles {
		namesMaxLen = max(namesMaxLen, len(profiles[i].Name))
		usersMaxLen = max(usersMaxLen, len(profiles[i].Data.User))
		projectsMaxLen = max(projectsMaxLen, len(profiles[i].Data.Project))
	}
	namesMaxLen = max(namesMaxLen, len("/name/"))
	usersMaxLen = max(usersMaxLen, len("/user/"))
	projectsMaxLen = max(projectsMaxLen, len("/project/"))
	{
		header := formLineProfilesDetails("/name/", "/user/", "/project/", namesMaxLen, usersMaxLen, projectsMaxLen)
		fmt.Println(header)
		//fmt.Println(strings.Repeat("_", len(header)))
	}
	for i := range profiles {
		_, err := fmt.Println(formLineProfilesDetails(profiles[i].Name, profiles[i].Data.User, profiles[i].Data.Project, namesMaxLen, usersMaxLen, projectsMaxLen))
		if err != nil {
			return err
		}
	}
	return nil
}

func formLineProfilesDetails(name, user, project string, widthName, widthUser, widthProject int) string {
	nameField := name + strings.Repeat(" ", widthName-len(name)) + "  |"
	userField := "  " + user + strings.Repeat(" ", widthUser-len(user)) + "  |"
	projectField := "  " + project + strings.Repeat(" ", widthProject-len(project))
	return nameField + userField + projectField
}

func printGetOutput(p profile.Profile) error {
	output := "profile:  " + p.Name + "\n" +
		"user   :  " + p.Data.User + "\n" +
		"project:  " + p.Data.Project
	_, err := fmt.Println(output)
	return err
}

func askToOverwrite(name string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	_, err := fmt.Printf("Profile %s already exist. Do you really want to overwrite it? [y/[n]]", name)
	if err != nil {
		return false, err
	}
	answer, _ := reader.ReadString('\n')
	// err != nil means input ended not with \n, thats acceptable

	answer = strings.TrimSpace(strings.ToLower(answer))

	return answer == "y" || answer == "yes", nil
}
