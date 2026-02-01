package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"cm-cli/pkg"
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
		name := fs.String("name", "default-cluster", "cluster name")
		root := fs.String("root", "configs", "root folder for configs")
		fmt.Println("Generating with name:", *name, "and root:", *root)
		// allow flags before or after filename by extracting filename first
		tokens := os.Args[2:]
		fnameIdx := -1
		for i, t := range tokens {
			if !strings.HasPrefix(t, "-") {
				fnameIdx = i
				break
			}
		}
		if fnameIdx == -1 {
			fmt.Println("usage: cm-cli generate <filename> [--name NAME] [--root ROOT]")
			os.Exit(1)
		}
		filename := tokens[fnameIdx]
		// build slice of flag tokens (everything except the filename)
		flagTokens := append([]string{}, tokens[:fnameIdx]...)
		if fnameIdx+1 < len(tokens) {
			flagTokens = append(flagTokens, tokens[fnameIdx+1:]...)
		}
		fs.Parse(flagTokens)

		if err := pkg.Generate(filename, *name, *root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "modify":
		fs := flag.NewFlagSet("modify", flag.ExitOnError)
		keyvalue := fs.String("keyvalue", "", "comma separated list of key=value pairs")
		root := fs.String("root", "configs", "root folder for configs")

		tokens := os.Args[2:]
		fnameIdx := -1
		for i, t := range tokens {
			if !strings.HasPrefix(t, "-") {
				fnameIdx = i
				break
			}
		}
		if fnameIdx == -1 {
			fmt.Println("usage: cm-cli modify <filename> --keyvalue key=val[,k2=v2] [--root ROOT]")
			os.Exit(1)
		}
		filename := tokens[fnameIdx]
		flagTokens := append([]string{}, tokens[:fnameIdx]...)
		if fnameIdx+1 < len(tokens) {
			flagTokens = append(flagTokens, tokens[fnameIdx+1:]...)
		}
		fs.Parse(flagTokens)

		if err := pkg.Modify(filename, *keyvalue, *root); err != nil {
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
	case "help":
		fmt.Println(pkg.Help())
	default:
		fmt.Println(pkg.Help())
		os.Exit(1)
	}
}
