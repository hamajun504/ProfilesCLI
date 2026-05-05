package cli

import "fmt"

func Run(args []string) error {
	if len(args) == 0 {
		if err := runProfile([]string{"help"}); err != nil {
			return err
		}
		return nil
	}
	if args[0] == "profile" {
		if err := runProfile(args[1:]); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("unknown command: %s", args[0])
}
