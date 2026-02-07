package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jbcool17/confighandler/cm-cli/pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println(pkg.Help())
		os.Exit(1)
	}

	cmd := os.Args[1]
	fmt.Println("Command:", cmd, "Args:", os.Args[2:])
	switch cmd {
	case "generate":
		fs := flag.NewFlagSet("generate", flag.ExitOnError)
		filename := fs.String("filename", "test.yaml", "output filename")
		name := fs.String("name", "default-cluster", "cluster name")
		root := fs.String("root", "configs", "root folder for configs")
		fmt.Println("Generating with name:", *name, "and root:", *root)

		if err := pkg.Generate(*filename, *name, *root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "modify":
		fs := flag.NewFlagSet("modify", flag.ExitOnError)
		filename := fs.String("filename", "test.yaml", "output filename")
		keyvalue := fs.String("keyvalue", "", "comma separated list of key=value pairs")
		root := fs.String("root", "configs", "root folder for configs")

		// process flags
		fs.Parse(os.Args[2:])

		if err := pkg.Modify(*filename, *keyvalue, *root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "schema":
		fs := flag.NewFlagSet("schema", flag.ExitOnError)
		output := fs.String("output", "schemas/cluster.schema.json", "output path for schema file")
		fs.Parse(os.Args[2:])

		if err := pkg.GenerateSchema(*output); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "debugyaml":
		fs := flag.NewFlagSet("debugyaml", flag.ExitOnError)
		filename := fs.String("filename", "test.yaml", "output filename")
		root := fs.String("root", "configs", "root folder for configs")

		if err := pkg.DebugYAML(*filename, *root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "help":
		fmt.Println(pkg.Help())
	default:
		fmt.Println(pkg.Help())
		os.Exit(1)
	}
}
