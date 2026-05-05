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
