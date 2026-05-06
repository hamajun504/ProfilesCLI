package profile

// ProfileData describes data stored in a profile YAML file.
type ProfileData struct {
	User    string `yaml:"user"`
	Project string `yaml:"project"`
}

// Profile describes a named profile with its YAML data.
type Profile struct {
	Name string
	Data ProfileData
}

func newProfile(name, user, project string) Profile {
	return Profile{
		Name: name,
		Data: ProfileData{
			User:    user,
			Project: project,
		},
	}
}
func newDefaultProfile(name string) Profile {
	return Profile{
		Name: name,
		Data: ProfileData{},
	}
}
