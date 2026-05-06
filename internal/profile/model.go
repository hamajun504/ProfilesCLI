package profile

type ProfileData struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

type Profile struct {
	Name string
	Data ProfileData
}

func NewProfile(name, user, project string) Profile {
	return Profile{
		Name: name,
		Data: ProfileData{
			User:    user,
			Project: project,
		},
	}
}
func NewDefaultProfile(name string) Profile {
	return Profile{
		Name: name,
		Data: ProfileData{},
	}
}
