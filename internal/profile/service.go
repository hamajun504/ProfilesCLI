package profile

import (
	"errors"
	"fmt"
	"os"
)

func Create(name, user, project string) error {
	if err := validateNewName(name); err != nil {
		return err
	}
	if err := validateUser(user); err != nil {
		return err
	}
	if err := validateProject(project); err != nil {
		return err
	}
	exist, err := Exist(name)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("profile %q already exists", name)
	}
	p := Profile{
		User:    user,
		Project: project,
	}
	if err := Save(name, p); err != nil {
		return err
	}
	return nil
}

func Update(name, user, project string) error {
	if err := validateNewName(name); err != nil {
		return err
	}
	if err := validateUser(user); err != nil {
		return err
	}
	if err := validateProject(project); err != nil {
		return err
	}
	exist, err := Exist(name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("profile %q not exists", name)
	}
	p := Profile{
		User:    user,
		Project: project,
	}
	if err := Save(name, p); err != nil {
		return err
	}
	return nil
}

func Get(name string) (string, string, error) {
	if err := validateOldName(name); err != nil {
		return "", "", err
	}
	p, err := Load(name)
	if err != nil {
		return "", "", err
	}
	return p.User, p.Project, nil
}

func List(mode FileStructure) ([]string, []string, []string, error) {
	profileNames, profiles, err := SearchAll(".", mode)
	if err != nil {
		return []string{}, []string{}, []string{}, err
	}
	users := make([]string, 0, len(profileNames))
	projects := make([]string, 0, len(profileNames))

	for i := range profileNames {
		users = append(users, profiles[i].User)
		projects = append(projects, profiles[i].Project)
	}

	return profileNames, users, projects, nil
}

func Delete(name string) error {
	if err := validateOldName(name); err != nil {
		return err
	}
	err := Remove(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	return nil
}
