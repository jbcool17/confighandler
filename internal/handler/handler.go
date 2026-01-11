package handler

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type TestConfig struct {
	FieldOne  int      `yaml:"fieldOne,omitempty"`
	Name      string   `yaml:"name"`
	Enabled   bool     `yaml:"enabled"`
	ListField []string `yaml:"listField,omitempty"`
	Options   options  `yaml:"options,omitempty"`
}

type TestConfigs struct {
	TestConfigs []TestConfig `yaml:"testConfigs"`
}

type options struct {
	Verbose bool
	Debug   bool
	Timeout int
	Version string
}

func Generate() {
	// 1. Create an instance of the struct and populate it
	configData := TestConfig{
		FieldOne:  42,
		Name:      "exampleConfig",
		Enabled:   true,
		ListField: []string{"item1", "item2", "item3"},
		Options: options{
			Verbose: true,
			Debug:   false,
			Timeout: 30,
			Version: "1.0.0",
		},
	}

	configs := TestConfigs{
		TestConfigs: []TestConfig{configData, configData, configData},
	}

	// 2. Marshal the struct into a YAML byte slice
	yamlData, err := yaml.Marshal(&configs)
	if err != nil {
		log.Fatalf("Error while Marshaling: %v", err)
	}

	// 3. Write the byte slice to a file
	fileName := "config.yaml"
	err = os.WriteFile(fileName, yamlData, 0644) // 0644 gives read/write permissions to owner, read to others
	if err != nil {
		log.Fatalf("Error writing to file %s: %v", fileName, err)
	}

	fmt.Printf("Successfully wrote YAML data to %s\n", fileName)
	fmt.Println("--- Generated YAML Content ---")
	fmt.Println(string(yamlData))
}
