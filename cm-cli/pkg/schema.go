package pkg

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schemas/cluster.schema.json
var clusterSchemaBytes []byte

// GetClusterSchema returns the JSON Schema loader for Cluster
func GetClusterSchema() gojsonschema.JSONLoader {
	return gojsonschema.NewBytesLoader(clusterSchemaBytes)
}

// ValidateClusterJSON validates a JSON object against the Cluster schema
func ValidateClusterJSON(data interface{}) error {
	// Convert data to JSON and back to ensure it's in the right format
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	documentLoader := gojsonschema.NewBytesLoader(jsonBytes)
	schemaLoader := GetClusterSchema()

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("schema validation error: %w", err)
	}

	if !result.Valid() {
		var errorMsg string
		for _, desc := range result.Errors() {
			errorMsg += fmt.Sprintf("  - %s\n", desc.String())
		}
		return fmt.Errorf("validation failed:\n%s", errorMsg)
	}

	return nil
}

// ValidateClusterYAML validates a YAML document (as map) against the Cluster schema
func ValidateClusterYAML(data map[string]interface{}) error {
	return ValidateClusterJSON(data)
}

// GenerateSchemaFromStruct generates a JSON Schema from the Cluster struct and writes it to a file
func GenerateSchemaFromStruct(outputPath string) error {
	schema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Cluster Configuration",
		"type":    "object",
	}

	properties := make(map[string]interface{})
	required := []string{}

	t := reflect.TypeOf(&Cluster{}).Elem()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")
		fieldName := strings.Split(yamlTag, ",")[0]
		if fieldName == "" {
			fieldName = field.Name
		}

		prop := generatePropertyFromType(field.Type)
		properties[fieldName] = prop
		required = append(required, fieldName)
	}

	schema["properties"] = properties
	schema["required"] = required

	// Marshal to JSON with pretty printing
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, append(schemaBytes, '\n'), 0644); err != nil {
		return fmt.Errorf("failed to write schema file: %w", err)
	}

	fmt.Printf("Generated schema file: %s\n", outputPath)
	return nil
}

func generatePropertyFromType(t reflect.Type) map[string]interface{} {
	prop := make(map[string]interface{})

	// Dereference pointers
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.String:
		prop["type"] = "string"
	case reflect.Bool:
		prop["type"] = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		prop["type"] = "integer"
	case reflect.Float32, reflect.Float64:
		prop["type"] = "number"
	case reflect.Slice, reflect.Array:
		prop["type"] = "array"
		elemType := t.Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() == reflect.Struct {
			items := make(map[string]interface{})
			items["type"] = "object"
			items["properties"] = generateStructProperties(elemType)
			items["required"] = getRequiredFields(elemType)
			prop["items"] = items
		}
	case reflect.Struct:
		prop["type"] = "object"
		prop["properties"] = generateStructProperties(t)
		prop["required"] = getRequiredFields(t)
	default:
		prop["type"] = "string"
	}

	return prop
}

func generateStructProperties(t reflect.Type) map[string]interface{} {
	properties := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")
		fieldName := strings.Split(yamlTag, ",")[0]
		if fieldName == "" || fieldName == "-" {
			continue
		}

		properties[fieldName] = generatePropertyFromType(field.Type)
	}

	return properties
}

func getRequiredFields(t reflect.Type) []string {
	required := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		yamlTag := field.Tag.Get("yaml")
		fieldName := strings.Split(yamlTag, ",")[0]
		if fieldName == "" || fieldName == "-" {
			continue
		}

		required = append(required, fieldName)
	}

	return required
}
