package env

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jbcool17/confighandler/utils"
)

type EnvConfigs struct {
	EnvConfigs []EnvConfig `yaml:"envConfigs"`
}

type EnvConfig struct {
	Name   string `yaml:"name"`
	Folder string `yaml:"folder"`
}

// CreateEnv writes an EnvConfig YAML file under the project's env folder.
// It returns the path to the created file.
func CreateEnv(name string) (string, error) {
	ec := EnvConfig{
		Name:   name,
		Folder: filepath.Join("configs", name),
	}

	data, err := yaml.Marshal(&ec)
	if err != nil {
		return "", err
	}

	// ensure env folder exists
	if err := os.MkdirAll(utils.EnvFolder, os.ModePerm); err != nil {
		return "", err
	}

	filePath := filepath.Join(utils.EnvFolder, name+".yaml")
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	return filePath, nil
}
