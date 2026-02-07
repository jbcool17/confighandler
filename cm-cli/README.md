# cm-cli

A simple configuration manager for YAML files with JSON schema validation.

## Features

- **Generate**: Create default cluster YAML configurations
- **Modify**: Update YAML fields using dot-notation with schema validation
- **Schema**: Auto-generate JSON schemas from Go structs
- **Validation**: Validate YAML against JSON schema with type checking

## Commands

### generate
```bash
go run . generate <filename> [--name NAME] [--root ROOT]
```
Creates a default cluster YAML file.

### modify
```bash
go run . modify <filename> --keyvalue key=value[,k2=v2] [--root ROOT]
```
Updates YAML fields with dot-notation support (e.g., `network_config.pod_cidr=10.0.0.0/16`).
Validates changes against schema before writing.

### schema
```bash
go run . schema [--output PATH]
```
Generates JSON schema from struct definitions. Updates schema when structs change.

### help
```bash
go run . help
```
Shows usage information.

## Examples

```bash
# Generate default config
go run . generate my-cluster.yaml --name my-cluster

# Modify single field
go run . modify my-cluster.yaml --keyvalue location=us-west1-a

# Modify multiple fields
go run . modify my-cluster.yaml --keyvalue "location=us-west1-a,feature_one.enabled=true"

# Update schema from structs
go run . schema
```

## Testing

```bash
go test ./cm-cli/pkg -v
```

## ToDo

- multiple files
