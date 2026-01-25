package oci

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jbcool17/confighandler/internal/tar"
	"github.com/joho/godotenv"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"gopkg.in/yaml.v3"
	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

func Execute() {
	fmt.Println("Pushing OCI artifact...")
	fileNames := []string{"test-demo.tar.gz", "test-directory"}
	tag := "0.3.0"
	push(fileNames, tag)
	pull(tag)
}

func connect() *remote.Repository {
	// Connect to a remote repository
	reg := os.Getenv("DOCKER_DOMAIN")
	repo, err := remote.NewRepository(reg + "/" + os.Getenv("DOCKER_USER") + "/demo-oci")
	if err != nil {
		panic(err)
	}

	repo.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.NewCache(),
		Credential: auth.StaticCredential(reg, auth.Credential{
			Username: os.Getenv("DOCKER_USER"),
			Password: os.Getenv("DOCKER_TOKEN"),
		}),
	}

	return repo
}

// func CopyOCI(ctx context.Context, src content.ReadOnlyStorage, dst content.Storage, ref v1.Descriptor) (ocispec.Descriptor, error) {

// 	log.Printf("Copying OCI artifact: %s\n", ref)

// 	err := oras.CopyGraph(
// 		ctx,
// 		src,
// 		dst,
// 		ref,
// 		oras.DefaultCopyGraphOptions,
// 	)
// 	if err != nil {
// 		return ocispec.Descriptor{}, err
// 	}
// 	log.Printf("Successfully copied OCI artifact: %s\n", ref)
// 	return ref, nil
// }

func copyToRemote(ctx context.Context, fs *file.Store, repo *remote.Repository, tag string) v1.Descriptor {
	// Copy from the file store to the remote repository
	log.Printf("Copying OCI artifact with tag: %s %s\n", tag, repo.Reference)

	// Copy from the file store to the remote repository
	manifestDescriptor, err := oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully copied OCI artifact with tag: %s %s\n", tag, repo.Reference)
	return manifestDescriptor
}

func copyFromRemote(ctx context.Context, fs *file.Store, repo *remote.Repository, tag string) v1.Descriptor {
	// Copy from the file store to the remote repository
	log.Printf("Copying OCI artifact with tag: %s %s\n", tag, repo.Reference)

	// Copy from the remote repository to the file store
	manifestDescriptor, err := oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	}
	log.Printf("Successfully copied OCI artifact with tag: %s %s\n", tag, repo.Reference)
	return manifestDescriptor
}

func push(fileNames []string, tag string) {
	// 0. Create a file store
	fs, err := file.New("tmp/")
	if err != nil {
		panic(err)
	}
	defer fs.Close()
	ctx := context.Background()

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 1. Add files to the file store
	mediaType := "application/vnd.test.file"
	// fileNames //:= manifestsPath //[]string{"test-demo.tar.gz", "tmp/myfile", "test-directory"}
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))
	for _, name := range fileNames {
		log.Printf("Adding file to OCI artifact: %s\n", name)
		fileDescriptor, err := fs.Add(ctx, name, mediaType, "")
		if err != nil {
			panic(err)
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("file descriptor for %s: %v\n", name, fileDescriptor)
	}

	log.Printf("Total files added to OCI artifact: %d\n", len(fileDescriptors))
	// log.Printf("%v", fileDescriptors)

	// 2. Pack the files and tag the packed manifest
	artifactType := "application/vnd.test.artifact"
	opts := oras.PackManifestOptions{
		Layers: fileDescriptors,
	}
	manifestDescriptor, err := oras.PackManifest(ctx, fs, oras.PackManifestVersion1_1, artifactType, opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("manifest descriptor:", manifestDescriptor)

	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		panic(err)
	}

	// 3. Connect to a remote repository
	repo := connect()
	// 4. Copy from the file store to the remote repository
	copyToRemote(ctx, fs, repo, tag)
}

func pull(tag string) {
	fmt.Println("Pulling OCI artifact...")
	// 0. Create a file store
	fs, err := file.New("tmp-out/")
	if err != nil {
		panic(err)
	}
	defer fs.Close()

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// 1. Connect to a remote repository
	ctx := context.Background()
	repo := connect()

	// 2. Copy from the remote repository to the file store
	manifestDescriptor := copyFromRemote(ctx, fs, repo, tag)
	fmt.Println("manifest descriptor:", manifestDescriptor)

	// 3. List the files pulled to the file store
	files, err := os.ReadDir("tmp-out/")
	if err != nil {
		log.Fatalf("Error reading env directory: %v", err)
	}

	for _, file := range files {
		fmt.Println("Pulled file:", file.Name())
	}

	// 4. (Optional) Read a tar file from the file store
	fileName := "test-demo.tar.gz"
	filePath := "tmp-out/" + fileName
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}
	fmt.Printf("Read file %s, size: %d bytes\n", fileName, len(data))
	tar.UntarGzFile(filePath, "tmp-out/untarred/")
	// get tmp-out/untarred/demo/test.yaml and import to env structs
	demoTestFile := "tmp-out/untarred/demo/test.yaml"
	data, err = os.ReadFile(demoTestFile)
	if err != nil {
		log.Fatalf("Error reading file %s: %v", demoTestFile, err)
	}

	var testConfig TestConfig
	err = yaml.Unmarshal(data, &testConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML from file %s: %v", demoTestFile, err)
	}

	fmt.Printf("Unmarshaled config from %s: %+v\n", demoTestFile, testConfig)

	// Remove tmp-out/untarred/ directory
	err = os.RemoveAll("tmp-out/untarred/")
	if err != nil {
		log.Fatalf("Error removing directory tmp-out/untarred/: %v", err)
	}

}

type TestConfig struct {
	Test string `yaml:"test"`
}
