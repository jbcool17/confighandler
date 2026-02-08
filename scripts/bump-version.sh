#!/usr/bin/env bash
# Script to bump version using semantic versioning
# Usage: ./scripts/bump-version.sh (patch|minor|major)

set -e

if [ $# -ne 1 ]; then
    echo "Usage: $0 (patch|minor|major)"
    exit 1
fi

TYPE=$1

# Get current version from the latest git tag (defaults to v0.4.0 if no tags exist)
CURRENT=$(git describe --tags --match='v*' --abbrev=0 2>/dev/null || echo "v0.4.0")
CURRENT=${CURRENT#v}

# Parse version
IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT"

case "$TYPE" in
    patch)
        PATCH=$((PATCH + 1))
        ;;
    minor)
        MINOR=$((MINOR + 1))
        PATCH=0
        ;;
    major)
        MAJOR=$((MAJOR + 1))
        MINOR=0
        PATCH=0
        ;;
    *)
        echo "Invalid type: $TYPE"
        exit 1
        ;;
esac

NEW_VERSION="v${MAJOR}.${MINOR}.${PATCH}"

echo "Bumping version: $CURRENT -> ${NEW_VERSION#v}"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
git push origin "$NEW_VERSION"

echo "Tagged and pushed $NEW_VERSION"
echo "GoReleaser will automatically create a release"
