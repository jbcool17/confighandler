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

# View version info
go run . version
```

## Version Management

Check the current version:
```bash
confighandler version
```

For detailed version management (bumping versions, releases, and build info), see [VERSION_MANAGEMENT.md](VERSION_MANAGEMENT.md).

### Releasing

Use the **bump-version** workflow in GitHub Actions:
1. Go to **Actions** → **bump-version**
2. Click **Run workflow**
3. Select `patch`, `minor`, or `major`
4. The workflow creates a new tag and triggers an automatic release

Or manually bump locally:
```bash
./scripts/bump-version.sh patch   # or minor, major
```

## Testing

```bash
go test ./cm-cli/pkg -v
```

## Install

```bash
go install github.com/jbcool17/confighandler@latest
```
