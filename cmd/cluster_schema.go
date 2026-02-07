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

// clusterSchemaCmd represents the clusterSchema command
var clusterSchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate cluster configuration schema",
	Long: `This command generates a JSON schema file for cluster configurations.
The schema defines the structure and validation rules for cluster configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if err := pkg.GenerateSchema(output); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterSchemaCmd)

	clusterSchemaCmd.Flags().String("output", "schemas/cluster.schema.json", "output path for schema file")
}
