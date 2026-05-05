package profile

import (
	"errors"
	"os"
	"strings"

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

var ErrNotYaml = errors.New("the file is not a yaml")

func SearchAll(path string) ([]string, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return []string{}, err
	}
	profileNames := make([]string, 0, len(dirEntry))
	for i := range dirEntry {
		name, err := getProfileName(dirEntry[i].Name())
		if errors.Is(err, ErrNotYaml) {
			continue
		}
		if err != nil {
			return profileNames, err
		}

		profileNames = append(profileNames, name)

	}
	return profileNames, nil
}

func getProfileName(fileName string) (string, error) {
	name, found := strings.CutSuffix(fileName, ".yaml")
	if found {
		return name, nil
	}
	name, found = strings.CutSuffix(fileName, ".yml")
	if found {
		return name, nil
	}
	return "", ErrNotYaml
}
