/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jbcool17/confighandler/cm-cli/pkg"
	"github.com/spf13/cobra"
)

// clusterGenCmd represents the clusterGen command
var clusterGenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate cluster configuration files",
	Long: `This command generates cluster configuration YAML files based on predefined templates and parameters.
It can be used to create default cluster configurations or modify existing ones.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Generating with name:", name, "and root:", root)

		if err := pkg.Generate(filename, name, root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterGenCmd)

	clusterGenCmd.Flags().String("filename", "test.yaml", "output filename")
	clusterGenCmd.Flags().String("name", "default-cluster", "cluster name")
	clusterGenCmd.Flags().String("root", "configs", "root folder for configs")
}
