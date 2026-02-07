# confighandler

A small CLI for managing YAML configuration files with JSON schema support.

## Features

- Generate default cluster YAMLs
- Modify YAML fields using dot-notation with schema validation
- Produce JSON schemas from Go structs
- Debug and inspect YAML files

## Commands

- `cluster gen` — generate a cluster YAML
- `cluster mod` — modify an existing cluster YAML
- `cluster schema` — output the cluster JSON schema
- `cluster debugyaml` — parse/debug a cluster YAML
- `env generate` — run environment-aware generation
- `env setup` — create an `env/<name>.yaml` entry (interactive or `-n` flag)

## Examples

```bash
go run . cluster gen --filename=my.yaml --name=my-cluster --root=configs
go run . cluster mod --filename=my.yaml --keyvalue "a=b" --root=configs
go run . cluster schema --output=schemas/cluster.schema.json
go run . cluster debugyaml --filename=my.yaml --root=configs

go run . env generate
go run . env setup --name=test
```

## Testing

```bash
go test ./cm-cli/pkg -v
```

## Install

```bash
go install github.com/jbcool17/confighandler@latest
```
