package handler

import (
	"fmt"
	"jbcool17/confighandler/internal/env"
	"jbcool17/confighandler/utils"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func GenerateHandler() {
	if doEnvConfigsExist() {
		fmt.Println("Environment configurations found. Proceeding with generation.")
		envConfigs := getEnvConfigs()
		fmt.Printf("Loaded environment configurations: %+v\n", envConfigs)
		for _, envConfig := range envConfigs.EnvConfigs {
			fmt.Printf("Generating config for environment: %s in folder: %s\n", envConfig.Name, envConfig.Folder)
			// Here you can add logic to generate different configs based on envConfig
			Generate(envConfig)
		}

	} else {
		fmt.Println("No environment configurations found. Generating default config.")
		Generate(env.EnvConfig{
			Name:   "default",
			Folder: "configs/default",
		})
	}
}

func Generate(envConfig env.EnvConfig) {
	// 1. Create an instance of the struct and populate it
	configData := testConfig{
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

	configs := testConfigs{
		TestConfigs: []testConfig{configData, configData, configData},
	}

	// 2. Marshal the struct into a YAML byte slice
	yamlData, err := yaml.Marshal(&configs)
	if err != nil {
		log.Fatalf("Error while Marshaling: %v", err)
	}

	// 3. Write the byte slice to a file
	fileName := "config.yaml"
	filePath := utils.RootDir + "/" + envConfig.Folder
	fullFilePath := filePath + "/" + fileName

	// Create the directory if it doesn't exist
	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating directory %s: %v", filePath, err)
	}

	err = os.WriteFile(fullFilePath, yamlData, 0644) // 0644 gives read/write permissions to owner, read to others
	if err != nil {
		log.Fatalf("Error writing to file %s: %v", fullFilePath, err)
	}

	fmt.Printf("Successfully wrote YAML data to %s\n", fullFilePath)
	// fmt.Println("--- Generated YAML Content ---")
	// fmt.Println(string(yamlData))
}

func doEnvConfigsExist() bool {
	// Check if "env" folder exists
	_, err := os.Stat(utils.EnvFolder)
	if os.IsNotExist(err) {
		fmt.Println("The 'env' directory does not exist.")
		return false
	}
	// Check if *.yaml files exists under the "env" folder
	files, err := os.ReadDir(utils.EnvFolder)
	if err != nil {
		log.Fatalf("Error reading env directory: %v", err)
	}

	if len(files) == 0 {
		fmt.Println("No YAML files found in the env directory.")
		return false
	}

	fmt.Println("YAML files in the env directory:")
	return true
}

func getEnvConfigs() env.EnvConfigs {
	var envConfigs env.EnvConfigs
	// Read all YAML files in the "env" folder
	files, err := os.ReadDir(utils.EnvFolder)
	if err != nil {
		log.Fatalf("Error reading env directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := utils.EnvFolder + "/" + file.Name()
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatalf("Error reading file %s: %v", filePath, err)
		}

		var singleEnvConfig env.EnvConfig
		err = yaml.Unmarshal(data, &singleEnvConfig)
		if err != nil {
			log.Fatalf("Error unmarshaling YAML from file %s: %v", filePath, err)
		}

		envConfigs.EnvConfigs = append(envConfigs.EnvConfigs, singleEnvConfig)
		fmt.Printf("Loaded config from %s: %+v\n", filePath, singleEnvConfig)
	}

	return envConfigs
}
