VERSION="v0.2.0"
git tag -a $VERSION -m "Add cluster configuration handler"
git push origin $VERSION


# if abovec succesfful then run goreleaser release
if [ $? -eq 0 ]; then
    goreleaser release --rm-dist
else
    echo "Git tagging failed. Aborting release."
fi