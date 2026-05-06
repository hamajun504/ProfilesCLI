package profile

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Save(name string, p ProfileData) error {
	data, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	nameFile := getFileName(name)
	if err = os.WriteFile(nameFile, data, 0o644); err != nil {
		return err
	}
	return nil
}

func Load(name string) (ProfileData, error) {
	nameFile := getFileName(name)
	data, err := os.ReadFile(nameFile)
	if err != nil {
		return ProfileData{}, err
	}
	p := ProfileData{}
	err = yaml.Unmarshal(data, &p)
	if err != nil {
		return ProfileData{}, err
	}
	return p, nil
}

var ErrNotYaml = errors.New("the file is not a yaml")

func SearchAll(path string, mode FileStructure) ([]string, []ProfileData, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return []string{}, []ProfileData{}, err
	}

	shouldInclude := getProfileFilter(mode)

	profileNames := make([]string, 0, len(dirEntry))
	profiles := make([]ProfileData, 0, len(dirEntry))

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
		if shouldInclude(struc) {
			profileNames = append(profileNames, name)
			profiles = append(profiles, prof)
		}
	}
	return profileNames, profiles, nil
}

func getProfileFilter(mode FileStructure) func(FileStructure) bool {
	switch mode {
	case Valid:
		return func(struc FileStructure) bool {
			return struc == Valid
		}

	case ValidOrExtended:
		return func(struc FileStructure) bool {
			return struc == Valid || struc == ValidOrExtended
		}

	default:
		return func(struc FileStructure) bool {
			return true
		}
	}
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
	Valid FileStructure = iota
	ValidOrExtended
	All
)

func validateFileStructure(path string) (FileStructure, ProfileData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return All, ProfileData{}, err
	}
	decoderKnown := yaml.NewDecoder(bytes.NewReader(data))
	decoderKnown.KnownFields(true)
	decoderUnknown := yaml.NewDecoder(bytes.NewReader(data))
	decoderUnknown.KnownFields(false)

	var p ProfileData

	if decoderUnknown.Decode(&p) != nil {
		return All, ProfileData{}, nil
	}
	if validateUser(p.User) != nil || validateProject(p.Project) != nil {
		return All, p, nil
	}
	pExtra := p
	if decoderKnown.Decode(&p) != nil {
		return ValidOrExtended, pExtra, nil
	}
	return Valid, p, nil
}
