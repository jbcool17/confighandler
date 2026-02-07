/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// envCmd represents the parent command for env related actions
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment related commands",
	Long:  `Group of commands to generate and manage environment configuration files.`,
}

func init() {
	rootCmd.AddCommand(envCmd)
}
