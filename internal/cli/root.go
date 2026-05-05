package cli

import "fmt"

func Run(args []string) error {
	if args[0] == "profile" {
		if err := runProfile(args[1:]); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Invalid command")
}
