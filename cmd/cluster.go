/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// clusterCmd represents the parent command for cluster related actions
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Cluster related commands",
	Long:  `Group of commands to generate, modify and inspect cluster configuration files.`,
}

func init() {
	rootCmd.AddCommand(clusterCmd)
}
