package profile

import (
	"bytes"
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
	if err = os.WriteFile(nameFile, data, 0644); err != nil {
		return err
	}
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

func SearchAllCorrect(path string) ([]string, error) {
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
		path, err = closePathWithSlash(path)
		if err != nil {
			return []string{}, err
		}
		struc, err := validateFileStructure(path + dirEntry[i].Name())
		if err != nil {
			return []string{}, err
		}
		if struc == Ok {
			profileNames = append(profileNames, name)
		}

	}
	return profileNames, nil
}

func SearchAllExtended(path string) ([]string, error) {
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
		path, err = closePathWithSlash(path)
		if err != nil {
			return []string{}, err
		}
		struc, err := validateFileStructure(path + dirEntry[i].Name())
		if err != nil {
			return []string{}, err
		}
		if struc == Ok || struc == ExtraFields {
			profileNames = append(profileNames, name)
		}
	}
	return profileNames, nil
}

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
		path, err = closePathWithSlash(path)
		if err != nil {
			return []string{}, err
		}
		profileNames = append(profileNames, name)
	}
	return profileNames, nil
}

func closePathWithSlash(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("attempt to access an empty path")
	}
	if path[len(path)-1] == '/' {
		return path, nil
	}
	return path + "/", nil
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
	fileName := getFileName(name)
	err := os.Remove(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintln(os.Stderr, "Profile not exist")
			return nil
		}
		return err
	}
	return nil

}

func Exist(name string) (bool, error) {
	_, err := os.Stat(getFileName(name))
	if err == nil {
		return true, nil
	}
	if err == os.ErrNotExist {
		return false, nil
	}
	return false, err

}

type FileStructure int

const (
	Ok FileStructure = iota
	Invalid
	ExtraFields
)

func validateFileStructure(path string) (FileStructure, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Invalid, err
	}
	decoder_kwn := yaml.NewDecoder(bytes.NewReader(data))
	decoder_kwn.KnownFields(true)
	decoder_unk := yaml.NewDecoder(bytes.NewReader(data))
	decoder_unk.KnownFields(false)

	var p Profile

	if decoder_unk.Decode(&p) != nil {
		return Invalid, nil
	}
	if validateUser(p.User) != nil || validateProject(p.Project) != nil {
		return Invalid, nil
	}
	if decoder_kwn.Decode(&p) != nil {
		return ExtraFields, nil
	}
	return Ok, nil
}
