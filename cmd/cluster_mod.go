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

// clusterModCmd represents the clusterMod command
var clusterModCmd = &cobra.Command{
	Use:   "mod",
	Short: "Modify existing cluster configuration files",
	Long: `This command modifies existing cluster configuration YAML files based on provided key-value pairs.
It allows users to update specific fields in the cluster configuration using dot notation for nested fields.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("filename")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		keyvalue, err := cmd.Flags().GetString("keyvalue")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		root, err := cmd.Flags().GetString("root")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if err := pkg.Modify(filename, keyvalue, root); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	clusterCmd.AddCommand(clusterModCmd)

	clusterModCmd.Flags().String("filename", "test.yaml", "output filename")
	clusterModCmd.Flags().String("keyvalue", "", "comma separated list of key=value pairs")
	clusterModCmd.Flags().String("root", "configs", "root folder for configs")
}
