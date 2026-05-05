package profile

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func SearchAllCorrect(path string) ([]string, []Profile, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return []string{}, []Profile{}, err
	}
	profileNames := make([]string, 0, len(dirEntry))
	profiles := make([]Profile, 0, len(dirEntry))
	for i := range dirEntry {
		name, err := getProfileName(dirEntry[i].Name())
		if err != nil {
			if errors.Is(err, ErrNotYaml) {
				continue
			}
			return profileNames, profiles, err
		}

		struc, prof, err := validateFileStructure(filepath.Join(path, dirEntry[i].Name()))
		if err != nil {
			return profileNames, profiles, err
		}
		if struc == Ok {
			profileNames = append(profileNames, name)
			profiles = append(profiles, prof)
		}

	}
	return profileNames, profiles, nil
}

func SearchAllExtended(path string) ([]string, []Profile, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return []string{}, []Profile{}, err
	}
	profileNames := make([]string, 0, len(dirEntry))
	profiles := make([]Profile, 0, len(dirEntry))
	for i := range dirEntry {
		name, err := getProfileName(dirEntry[i].Name())
		if err != nil {
			if errors.Is(err, ErrNotYaml) {
				continue
			}
			return profileNames, profiles, err
		}
		struc, prof, err := validateFileStructure(filepath.Join(path, dirEntry[i].Name()))
		if err != nil {
			return profileNames, profiles, err
		}
		if struc == Ok || struc == ExtraFields {
			profileNames = append(profileNames, name)
			profiles = append(profiles, prof)
		}
	}
	return profileNames, profiles, err
}

func SearchAll(path string) ([]string, []Profile, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return []string{}, []Profile{}, err
	}
	profileNames := make([]string, 0, len(dirEntry))
	profiles := make([]Profile, 0, len(dirEntry))
	for i := range dirEntry {
		name, err := getProfileName(dirEntry[i].Name())
		if err != nil {
			if errors.Is(err, ErrNotYaml) {
				continue
			}
			return profileNames, profiles, err
		}
		_, prof, err := validateFileStructure(filepath.Join(path, dirEntry[i].Name()))
		if err != nil {
			return profileNames, profiles, err
		}
		profiles = append(profiles, prof)
		profileNames = append(profileNames, name)
	}
	return profileNames, profiles, nil
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
	if errors.Is(err, os.ErrNotExist) {
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

func validateFileStructure(path string) (FileStructure, Profile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Invalid, Profile{}, err
	}
	decoder_kwn := yaml.NewDecoder(bytes.NewReader(data))
	decoder_kwn.KnownFields(true)
	decoder_unk := yaml.NewDecoder(bytes.NewReader(data))
	decoder_unk.KnownFields(false)

	var p Profile

	if decoder_unk.Decode(&p) != nil {
		return Invalid, Profile{}, nil
	}
	if validateUser(p.User) != nil || validateProject(p.Project) != nil {
		return Invalid, p, nil
	}
	p_ext := p
	if decoder_kwn.Decode(&p) != nil {
		return ExtraFields, p_ext, nil
	}
	return Ok, p, nil
}
