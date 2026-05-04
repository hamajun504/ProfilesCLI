package cli

func Run(args []string) {
	if args[0] == "profile" {
		runProfile(args[1:])
	}

}
