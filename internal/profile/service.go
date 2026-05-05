package profile

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

func List(flag FileStructure) ([]string, []string, []string, error) {
	var err error
	var profileNames []string
	var profiles []Profile
	switch flag {
	case Ok:
		profileNames, profiles, err = SearchAllCorrect(".")
	case ExtraFields:
		profileNames, profiles, err = SearchAllExtended(".")
	default:
		profileNames, profiles, err = SearchAll(".")
	}
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
		return err
	}
	return nil
}
