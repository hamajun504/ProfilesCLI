package profile

func Create(name, user, project string) error {
	if err := validateNewName(name); err != nil {
		return err
	}
	if err := validateUser(name); err != nil {
		return err
	}
	if err := validateProject(name); err != nil {
		return err
	}

	p := Profile{
		User:    user,
		Project: project,
	}
	Save(name, p)
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

func List() ([]string, error) {
	profileNames, err := SearchAll(".")
	if err != nil {
		return []string{}, err
	}
	return profileNames, nil
}

func Delete(name string) error {
	err := Remove(name)
	if err != nil {
		return err
	}
	return nil
}
