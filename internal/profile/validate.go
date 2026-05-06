package profile

import (
	"fmt"
	"regexp"
	"strings"
)

var profileNameRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,64}$`)

func validateNewName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if !profileNameRe.MatchString(name) {
		return fmt.Errorf("profile name must contain only letters, digits, '-' or '_' and be 1-64 characters long")
	}

	return nil
}

func validateOldName(name string) error {
	if name == "" {
		return fmt.Errorf("name is required")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("name must not contain path separators")
	}

	return nil
}

func validateUser(user string) error {
	user = strings.TrimSpace(user)

	if user == "" {
		return fmt.Errorf("user is required")
	}

	if strings.ContainsAny(user, "\r\n") {
		return fmt.Errorf("user must not contain line breaks")
	}

	return nil
}

func validateProject(project string) error {
	project = strings.TrimSpace(project)

	if project == "" {
		return fmt.Errorf("project is required")
	}

	if strings.ContainsAny(project, "\r\n") {
		return fmt.Errorf("project must not contain line breaks")
	}

	return nil
}
