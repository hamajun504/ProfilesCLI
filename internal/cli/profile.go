package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/hamajun504/ProfilesCLI/internal/profile"
)

func runProfile(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	switch args[0] {
	case "create":
		return runProfileCreate(args[1:])
	case "get":
		return runProfileGet(args[1:])
	case "list":
		return runProfileList(args[1:])
	case "delete":
		return runProfileDelete(args[1:])
	case "help":
		return printHelp()
	default:
		return fmt.Errorf("unknown profile command: %s", args[0])
	}
}

func runProfileCreate(args []string) error {
	fs := flag.NewFlagSet("profile create", flag.ContinueOnError)

	name := fs.String("name", "", "profile name")
	user := fs.String("user", "", "user name")
	project := fs.String("project", "", "project name")
	forceOverwrite := false
	fs.BoolVar(&forceOverwrite, "force", false, "overwrite existing profile")
	fs.BoolVar(&forceOverwrite, "f", false, "overwrite existing profile")

	if err := fs.Parse(args); err != nil {
		return err
	}
	exist, err := profile.Exists(*name)
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
	return nil
}

func runProfileGet(args []string) error {
	fs := flag.NewFlagSet("profile get", flag.ContinueOnError)
	name := fs.String("name", "", "profile name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	p, err := profile.Get(*name)
	if err != nil {
		return err
	}
	if err := printGetOutput(p); err != nil {
		return err
	}
	return nil
}

func runProfileList(args []string) error {
	fs := flag.NewFlagSet("profile list", flag.ContinueOnError)
	extendedFiles := fs.Bool("e", false, "output files with extra fields")
	allFiles := fs.Bool("a", false, "output all yaml-files")
	longOutput := fs.Bool("l", false, "detailed output")

	if err := fs.Parse(args); err != nil {
		return err
	}

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
	return nil
}

func runProfileDelete(args []string) error {
	fs := flag.NewFlagSet("profile delete", flag.ContinueOnError)
	name := fs.String("name", "", "profile name")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := profile.Delete(*name); err != nil {
		return err
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
		namesMaxLen = max(namesMaxLen, utf8.RuneCountInString(profiles[i].Name))
		usersMaxLen = max(usersMaxLen, utf8.RuneCountInString(profiles[i].Data.User))
		projectsMaxLen = max(projectsMaxLen, utf8.RuneCountInString(profiles[i].Data.Project))
	}
	namesMaxLen = max(namesMaxLen, utf8.RuneCountInString("/name/"))
	usersMaxLen = max(usersMaxLen, utf8.RuneCountInString("/user/"))
	projectsMaxLen = max(projectsMaxLen, utf8.RuneCountInString("/project/"))
	{
		header := formLineProfilesDetails("/name/", "/user/", "/project/", namesMaxLen, usersMaxLen, projectsMaxLen)
		if _, err := fmt.Println(header); err != nil {
			return err
		}
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
	nameField := name + strings.Repeat(" ", widthName-utf8.RuneCountInString(name)) + "  |"
	userField := "  " + user + strings.Repeat(" ", widthUser-utf8.RuneCountInString(user)) + "  |"
	projectField := "  " + project + strings.Repeat(" ", widthProject-utf8.RuneCountInString(project))
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
	return askToOverwriteFromReader(name, os.Stdin)
}

func askToOverwriteFromReader(name string, r io.Reader) (bool, error) {
	reader := bufio.NewReader(r)
	_, err := fmt.Printf("Profile %s already exists. Do you really want to overwrite it? [y/N] ", name)
	if err != nil {
		return false, err
	}
	answer, err := reader.ReadString('\n')
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return false, err
		}
	}

	answer = strings.TrimSpace(strings.ToLower(answer))

	return answer == "y" || answer == "yes", nil
}
