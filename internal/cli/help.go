package cli

import "fmt"

func printHelp() error {
	_, err := fmt.Println(helpMessage)
	return err
}

var helpMessage string = `mws profile - CLI for profiles managing
Usage:
  mws <command> [arguments]

Commands:
  profile create   Create a new profile
  profile get      Show profile by name
  profile list     List profiles in the current directory
  profile delete   Delete profile by name
  profile help     Show this help message
  help             Show this help message

Profile commands:
  mws profile create --name=<name> --user=<user> --project=<project> [--force|-f]
  mws profile get --name=<name>
  mws profile list [-a|-e|-l]
  mws profile delete --name=<name>
  mws profile help

Options:
  --name       Profile name. Used as the YAML file name without extension
  --user       User value saved into the profile
  --project    Project value saved into the profile
  --force, -f  Overwrite an existing profile without confirmation

List options:
  -l           Show only valid profiles
  -e           Show valid profiles and profiles with extra YAML fields
  -a           Show all YAML profile files, including invalid ones

Examples:
  mws profile create --name=test --user=example --project=new-project
  mws profile create --name=test --user=example --project=new-project --force
  mws profile get --name=test
  mws profile list
  mws profile list -l
  mws profile delete --name=test

Notes:
  Profiles are stored as YAML files in the current directory.
  The profile name corresponds to the file name: <name>.yaml.
  Each profile file contains two string fields: user and project.`
