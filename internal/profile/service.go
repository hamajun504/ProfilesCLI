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

func List(flag FileStructure) ([]string, error) {
	var err error
	var profileNames []string
	switch flag {
	case Ok:
		profileNames, err = SearchAllCorrect(".")
	case ExtraFields:
		profileNames, err = SearchAllExtended(".")
	default:
		profileNames, err = SearchAll(".")
	}
	if err != nil {
		return []string{}, err
	}
	return profileNames, nil
}

func Delete(name string) error {
	if err := validateOldName(name); err != nil {
		return err
	}
	err := Remove(name)
	if err != nil {
		return err
	}
	return nil
}
