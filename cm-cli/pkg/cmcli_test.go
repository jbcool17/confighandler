package pkg

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateCreatesFile(t *testing.T) {
	tmp := t.TempDir()
	filename := "test-cluster.yaml"
	if err := Generate(filename, "my-cluster", tmp); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	path := filepath.Join(tmp, filename)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file created, stat error: %v", err)
	}
}

func TestModifyUpdatesField(t *testing.T) {
	tmp := t.TempDir()
	filename := "cluster.yaml"
	if err := Generate(filename, "start-cluster", tmp); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	// modify location and nested network_config.pod_cidr
	kv := "location=us-west1-a,network_config.pod_cidr=10.1.0.0/16"
	if err := Modify(filename, kv, tmp); err != nil {
		t.Fatalf("Modify failed: %v", err)
	}
	// read file and ensure new values present
	b, err := os.ReadFile(filepath.Join(tmp, filename))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, "us-west1-a") {
		t.Fatalf("expected modified location in file")
	}
	if !strings.Contains(s, "10.1.0.0/16") {
		t.Fatalf("expected modified pod cidr in file")
	}
}

func TestModifyNodePoolIndexed(t *testing.T) {
	tmp := t.TempDir()
	filename := "cluster2.yaml"
	if err := Generate(filename, "start-cluster", tmp); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	// set node_pools[0].name
	kv := "node_pools[0].name=pool-zero,node_pools[0].machine_type=n1-standard-1"
	if err := Modify(filename, kv, tmp); err != nil {
		t.Fatalf("Modify failed: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(tmp, filename))
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	s := string(b)
	if !strings.Contains(s, "pool-zero") {
		t.Fatalf("expected modified node pool name in file")
	}
	if !strings.Contains(s, "n1-standard-1") {
		t.Fatalf("expected modified machine type in file")
	}
}

func TestValidateClusterJSONSchema(t *testing.T) {
	tests := []struct {
		name      string
		data      map[string]interface{}
		shouldErr bool
	}{
		{
			name: "valid cluster",
			data: map[string]interface{}{
				"name":     "test-cluster",
				"location": "us-central1-a",
				"network_config": map[string]interface{}{
					"network":  "default",
					"pod_cidr": "10.0.0.0/16",
				},
				"node_pools": []map[string]interface{}{
					{
						"name":         "pool-1",
						"machine_type": "e2-medium",
					},
				},
				"feature_one": map[string]interface{}{
					"enabled": true,
				},
			},
			shouldErr: false,
		},
		{
			name: "missing required field",
			data: map[string]interface{}{
				"location": "us-central1-a",
				"network_config": map[string]interface{}{
					"network":  "default",
					"pod_cidr": "10.0.0.0/16",
				},
			},
			shouldErr: true,
		},
		{
			name: "invalid pod_cidr format",
			data: map[string]interface{}{
				"name":     "test-cluster",
				"location": "us-central1-a",
				"network_config": map[string]interface{}{
					"network":  "default",
					"pod_cidr": "invalid-cidr",
				},
			},
			shouldErr: true,
		},
		{
			name: "wrong type for enabled field",
			data: map[string]interface{}{
				"name":     "test-cluster",
				"location": "us-central1-a",
				"network_config": map[string]interface{}{
					"network":  "default",
					"pod_cidr": "10.0.0.0/16",
				},
				"feature_one": map[string]interface{}{
					"enabled": "true", // should be boolean, not string
				},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClusterJSON(tt.data)
			if tt.shouldErr && err == nil {
				t.Fatalf("expected validation error, got none")
			}
			if !tt.shouldErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestGenerateSchemaFromStruct(t *testing.T) {
	tmp := t.TempDir()
	schemaPath := filepath.Join(tmp, "test-schema.json")

	if err := GenerateSchemaFromStruct(schemaPath); err != nil {
		t.Fatalf("GenerateSchemaFromStruct failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(schemaPath); err != nil {
		t.Fatalf("schema file not created: %v", err)
	}

	// Verify it's valid JSON
	b, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read schema file: %v", err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(b, &schema); err != nil {
		t.Fatalf("schema is not valid JSON: %v", err)
	}

	// Verify key schema properties
	if schema["$schema"] == nil {
		t.Fatal("schema missing $schema field")
	}
	if schema["title"] == nil {
		t.Fatal("schema missing title field")
	}
	if schema["type"] == nil {
		t.Fatal("schema missing type field")
	}
	if schema["properties"] == nil {
		t.Fatal("schema missing properties field")
	}
}
