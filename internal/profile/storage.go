package profile

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Save(name string, p Profile) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	nameFile := getFileName(name)
	os.WriteFile(nameFile, data, 0644)
	return nil
}

func Load(name string) (Profile, error) {
	nameFile := getFileName(name)
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
		if err != nil {
			if errors.Is(err, ErrNotYaml) {
				continue
			}
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
	return "", ErrNotYaml
}

func getFileName(name string) string {
	return name + ".yaml"
}

func Remove(name string) error {
	fileName := name + ".yaml"
	err := os.Remove(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "Profile not exist")
			return nil
		}
		if errors.Is(err, os.ErrPermission) {
			return fmt.Errorf("No permission for file deletion")
		}
	}
	return nil

}
