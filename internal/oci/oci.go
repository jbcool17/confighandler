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
	push()
	pull()

}
func push() {
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
	fileNames := []string{"test-demo.tar.gz", "tmp/myfile", "test-directory"}
	fileDescriptors := make([]v1.Descriptor, 0, len(fileNames))
	for _, name := range fileNames {
		fileDescriptor, err := fs.Add(ctx, name, mediaType, "")
		if err != nil {
			panic(err)
		}
		fileDescriptors = append(fileDescriptors, fileDescriptor)
		fmt.Printf("file descriptor for %s: %v\n", name, fileDescriptor)
	}

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

	tag := "0.1.0"
	if err = fs.Tag(ctx, manifestDescriptor, tag); err != nil {
		panic(err)
	}

	// 3. Connect to a remote repository
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

	// 4. Copy from the file store to the remote repository
	_, err = oras.Copy(ctx, fs, tag, repo, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	}
}

func pull() {
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

	// 2. Copy from the remote repository to the file store
	tag := "0.1.0"
	manifestDescriptor, err := oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
	if err != nil {
		panic(err)
	}
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
