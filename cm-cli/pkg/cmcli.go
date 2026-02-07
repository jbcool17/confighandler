package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Generate creates a default cluster yaml file in root/<filename>
func Generate(filename, name, root string) error {
	fmt.Println("Generating", filename, "with name", name, "in", root)
	if err := ensureRootDir(root); err != nil {
		return err
	}
	cluster := Cluster{
		Name:     name,
		Location: "us-central1-a",
		NetworkConfig: NetworkConfig{
			Network: "default",
			PodCIDR: "10.0.0.0/16",
		},
		NodePools: []NodePool{{Name: "default-pool", MachineType: "e2-medium"}},
	}
	target := filepath.Join(root, filename)
	return WriteYAMLFile(target, cluster)
}

// Modify updates YAML fields based on comma-separated key=value pairs.
// Keys use dot-notation for nested fields, for example network_config.pod_cidr=pod-cidr-2
// After modification, validates the document against JSON Schema.
func Modify(filename, keyvalue, root string) error {
	fmt.Println("Modifying", filename, "with", keyvalue, "in", root)
	if err := ensureRootDir(root); err != nil {
		return err
	}
	target := filepath.Join(root, filename)
	b, err := os.ReadFile(target)
	if err != nil {
		return fmt.Errorf("read error: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return fmt.Errorf("yaml unmarshal: %w", err)
	}
	if len(doc.Content) == 0 {
		return fmt.Errorf("empty document")
	}
	rootNode := doc.Content[0]

	pairs := parseKeyValuePairs(keyvalue)
	for k, v := range pairs {
		path := strings.Split(k, ".")
		if err := setYAMLPath(rootNode, path, v); err != nil {
			return fmt.Errorf("set path %s: %w", k, err)
		}
	}

	// Validate against JSON Schema before writing
	var data map[string]interface{}
	updatedBytes, err := yaml.Marshal(rootNode)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}
	if err := yaml.Unmarshal(updatedBytes, &data); err != nil {
		return fmt.Errorf("unmarshal for validation: %w", err)
	}

	if err := ValidateClusterJSON(data); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// write back
	return WriteYAMLFile(target, &doc)
}

// GenerateSchema generates a JSON schema file from the Cluster struct
func GenerateSchema(outputPath string) error {
	if outputPath == "" {
		outputPath = "schemas/cluster.schema.json"
	}
	return GenerateSchemaFromStruct(outputPath)
}

// Debug helper for YAML
func DebugYAML(filename string, root string) error {
	// root := "configs"
	// filename := "test-01.yaml"
	target := filepath.Join(root, filename)
	content, err := os.ReadFile(target)
	if err != nil {
		return fmt.Errorf("error reading file: %v\n", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return fmt.Errorf("yaml unmarshal: %w", err)
	}
	if len(doc.Content) == 0 {
		return fmt.Errorf("empty document")

	}
	fmt.Println("YAML Document Content Length:", len(doc.Content))
	fmt.Printf("YAML Document Content: %+v\n", doc)
	rootNode := doc.Content[0]

	fmt.Printf("Original YAML Node:\n%+v\n", rootNode)
	fmt.Println("-----")
	DumpYAMLNode(&doc)
	return nil
}

// Help returns the help message for cm-cli
func Help() string {
	return `cm-cli - simple config manager

Usage:
  cm-cli generate [--filename FILENAME] [--name NAME] [--root ROOT]
    Generate a default cluster YAML file. Default name: "default-cluster", default root: "configs".

  cm-cli modify [--filename FILENAME] [--keyvalue KEYVALUE] [--root ROOT]
    Modify fields in an existing YAML using comma-separated key=value pairs. Use dot notation for nested fields: network_config.pod_cidr=pod-cidr-2

  cm-cli schema [--output OUTPUT]
    Generate JSON schema from struct definitions. Default output: schemas/cluster.schema.json

  cm-cli debugyaml [--filename FILENAME] [--root ROOT]
	Debug and dump the YAML node structure of the specified file in the configs/ directory.

  cm-cli help
    Print this help message.
`
}
