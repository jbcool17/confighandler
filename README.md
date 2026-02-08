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

Use the **prepare-release** workflow in GitHub Actions:
1. Go to **Actions** → **prepare-release** → **Run workflow**
2. Select version type: `patch`, `minor`, or `major`
3. Toggle **trigger_release** to auto-release (default: on)
4. Workflow creates tag and optionally triggers GoReleaser

Or manually bump locally:
```bash
./scripts/bump-version.sh patch   # or minor, major
```

## Testing

```bash
go test ./cm-cli/pkg -v
```

## Install

**Recommended: Download pre-built binary**

```bash
VERSION=v0.4.5
OS=Darwin          # or Linux
ARCH=arm64         # or x86_64, i386, etc.

curl -L https://github.com/jbcool17/confighandler/releases/download/$VERSION/confighandler_${OS}_${ARCH}.tar.gz | tar xz
sudo mv confighandler /usr/local/bin/
```

See [GitHub Releases](https://github.com/jbcool17/confighandler/releases) for available versions and platform options.

**Alternative: Install from source**

```bash
go install github.com/jbcool17/confighandler@latest
```

Note: Installing from source will show basic version info only. For full version details, use pre-built binaries.
