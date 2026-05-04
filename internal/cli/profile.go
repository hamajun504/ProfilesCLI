package cli

import "fmt"

func runProfile(args []string) {
	switch args[0] {
	case "create":
		fmt.Println("Create Profile")
	case "get":
		fmt.Println("Return Profile")
	case "list":
		fmt.Println("List all profiles")
	case "delete":
		fmt.Println("Delete Profile")
	}
}
