package profile

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Save(name string, p Profile) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	nameFile := name + ".yaml"
	os.WriteFile(nameFile, data, 0644)
	return nil
}

func Load(name string) (Profile, error) {
	nameFile := name + ".yaml"
	data, err := os.ReadFile(nameFile)
	if err != nil {
		return Profile{}, err
	}
	p := Profile{}
	err = yaml.Unmarshal(data, &p)
	return p, nil
}
