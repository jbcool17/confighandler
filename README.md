# confighandler

- an example tool for managing and generating yaml

## Usage

- https://cobra.dev/

```
go run . generate-env -n test
go run . generate
```


## Release

- https://goreleaser.com/quick-start/

```
go install github.com/jbcool17/confighandler@latest
```


## OCI

```bash
tar -czf demo.tar.gz demo

oras login
oras push docker.io/jbcool17/myartifact:1.0.0 demo.tar.gz:application/vnd.oci.image.layer.v1.tar+gzip
```
