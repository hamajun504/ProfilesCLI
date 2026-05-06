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

func Save(p Profile) error {
	data, err := yaml.Marshal(p.Data)
	if err != nil {
		return err
	}
	nameFile := getFileName(p.Name)
	if err = os.WriteFile(nameFile, data, 0o644); err != nil {
		return err
	}
	return nil
}

func Load(name string) (Profile, error) {
	data, err := os.ReadFile(getFileName(name))
	if err != nil {
		return Profile{}, err
	}
	p := newDefaultProfile(name)
	err = yaml.Unmarshal(data, &p.Data)
	if err != nil {
		return Profile{}, err
	}
	return p, nil
}

var ErrNotYaml = errors.New("the file is not a yaml")

func SearchAll(path string, mode FileStructure) ([]Profile, error) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	shouldInclude := getProfileFilter(mode)

	profiles := make([]Profile, 0, len(dirEntry))

	for i := range dirEntry {
		name, err := getProfileName(dirEntry[i].Name())
		if err != nil {
			if errors.Is(err, ErrNotYaml) {
				continue
			}
			return profiles, err
		}
		prof := newDefaultProfile(name)
		struc, err := validateFileStructure(filepath.Join(path, dirEntry[i].Name()), &prof)
		if err != nil {
			return profiles, err
		}
		if shouldInclude(struc) {
			profiles = append(profiles, prof)
		}
	}
	return profiles, nil
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

func validateFileStructure(path string, p *Profile) (FileStructure, error) {
	if p == nil {
		return All, fmt.Errorf("profile is nil")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return All, err
	}
	decoderStrict := yaml.NewDecoder(bytes.NewReader(data))
	decoderStrict.KnownFields(true)
	decoderNotStrict := yaml.NewDecoder(bytes.NewReader(data))
	decoderNotStrict.KnownFields(false)

	if decoderNotStrict.Decode(&p.Data) != nil {
		return All, nil
	}
	if validateUser(p.Data.User) != nil || validateProject(p.Data.Project) != nil {
		return All, nil
	}
	pBeforeStrictDecode := *p
	if decoderStrict.Decode(&p.Data) != nil {
		*p = pBeforeStrictDecode
		return ValidOrExtended, nil
	}
	return Valid, nil
}
