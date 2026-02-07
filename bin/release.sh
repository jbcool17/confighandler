VERSION="v0.2.0"
git tag -a $VERSION -m "Add cluster configuration handler"
git push origin $VERSION


if [ $? -eq 0 ]; then
    goreleaser release
else
    echo "Git tagging failed. Aborting release."
fi