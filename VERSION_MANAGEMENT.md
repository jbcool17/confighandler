# Version Management

This document explains how versioning works in confighandler.

## Displaying Version

Run the CLI with the `version` command:

```bash
go run . version
```

In development, this shows:
```
confighandler dev (built unknown, commit unknown)
```

After building with GoReleaser (or using ldflags), it shows:
```
confighandler v1.2.3 (built 2026-02-08, commit abc1234)
```

## Building with Version Info

### Development Build
```bash
go build -o confighandler
./confighandler version
```

### Release Build (with ldflags)
```bash
go build \
  -ldflags="-X github.com/jbcool17/confighandler/internal/version.Version=v1.2.3 \
            -X github.com/jbcool17/confighandler/internal/version.BuildTime=$(date) \
            -X github.com/jbcool17/confighandler/internal/version.GitCommit=$(git rev-parse HEAD)" \
  -o confighandler
```

## Releasing & Version Bumping

### Automatic Release (GitHub Actions)

1. **Bump version** using the script:
   ```bash
   ./scripts/bump-version.sh patch   # or minor, major
   ```

   This creates and pushes a git tag like `v1.2.4`.

2. **GitHub Actions automatically**:
   - Detects the tag via the `goreleaser.yaml` workflow
   - Runs GoReleaser with the new version
   - Creates a GitHub Release with binaries
   - Injects version info via ldflags

### Manual Release

Create and push a git tag directly:
```bash
git tag -a v1.2.3 -m "Release v1.2.3"
git push origin v1.2.3
```

The goreleaser workflow will automatically trigger.

## Version Info Variables

Set via `-ldflags` during build:

- `github.com/jbcool17/confighandler/internal/version.Version` — release tag (e.g., `v1.2.3`)
- `github.com/jbcool17/confighandler/internal/version.BuildTime` — build timestamp
- `github.com/jbcool17/confighandler/internal/version.GitCommit` — commit hash

GoReleaser interpolates these automatically:
- `{{.Version}}` → matched git tag
- `{{.Date}}` → build timestamp
- `{{.Commit}}` → commit hash
