package profile

import (
	"errors"
	"fmt"
	"os"
	"sort"
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
	exist, err := Exists(name)
	if err != nil {
		return err
	}
	if exist {
		return fmt.Errorf("profile %q already exists", name)
	}
	p := newProfile(name, user, project)
	if err := Save(p); err != nil {
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
	exist, err := Exists(name)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("profile %q does not exists", name)
	}
	p := newProfile(name, user, project)
	if err := Save(p); err != nil {
		return err
	}
	return nil
}

func Get(name string) (Profile, error) {
	if err := validateOldName(name); err != nil {
		return Profile{}, err
	}
	p, err := Load(name)
	if err != nil {
		return Profile{}, err
	}
	return p, nil
}

func List(mode FileStructure) ([]Profile, error) {
	profiles, err := SearchAll(".", mode)
	if err != nil {
		return nil, err
	}

	sort.Slice(profiles, func(i, j int) bool {
		return profiles[i].Name < profiles[j].Name
	})

	return profiles, nil
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
