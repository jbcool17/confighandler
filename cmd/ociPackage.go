/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/jbcool17/confighandler/internal/oci"
	"github.com/jbcool17/confighandler/internal/tar"
	"github.com/spf13/cobra"
)

// ociPackageCmd represents the ociPackage command
var ociPackageCmd = &cobra.Command{
	Use:   "ociPackage",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ociPackage called")
		tar.ExecuteTar()
		oci.Execute()
	},
}

func init() {
	rootCmd.AddCommand(ociPackageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ociPackageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ociPackageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
