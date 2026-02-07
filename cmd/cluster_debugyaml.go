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

// clusterDebugYamlCmd represents the clusterDebugYaml command
var clusterDebugYamlCmd = &cobra.Command{
	Use:   "debugyaml",
	Short: "Debug cluster YAML configuration files",
	Long: `This command debugs cluster YAML configuration files by parsing and displaying their structure.
It helps identify YAML syntax errors and validates the configuration against the schema.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if err := pkg.DebugYAML(filename, root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterDebugYamlCmd)

	clusterDebugYamlCmd.Flags().String("filename", "test.yaml", "output filename")
	clusterDebugYamlCmd.Flags().String("root", "configs", "root folder for configs")
}
