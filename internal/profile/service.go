package profile

func Create(name, user, project string) {
	p := Profile{
		User:    user,
		Project: project,
	}
	Save(name, p)
}
