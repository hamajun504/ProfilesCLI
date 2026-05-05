package cli

var helpMessage string = `mws profile - CLI for profiles managing
Usage:
  mws profile create --name=<name> --user=<user> --project=<project>
  mws profile get --name=<name>
  mws profile list
  mws profile delete --name=<name>
  mws help

Commands:
  profile create   Create a profile
  profile get      Show profile by name
  profile list     List all profiles
  profile delete   Delete profile by name
  help             Show this help message`
