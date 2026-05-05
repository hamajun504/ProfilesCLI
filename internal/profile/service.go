package profile

func Create(name, user, project string) {
	p := Profile{
		User:    user,
		Project: project,
	}
	Save(name, p)
}

func Get(name string) (string, string, error) {
	p, err := Load(name)
	if err != nil {
		return "", "", err
	}
	return p.User, p.Project, nil
}
