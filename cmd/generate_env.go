package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"jbcool17/confighandler/internal/env"

	"github.com/spf13/cobra"
)

// generateEnvCmd represents the generate-env command
var generateEnvCmd = &cobra.Command{
	Use:   "generate-env [name]",
	Short: "Create an env YAML file for an environment",
	Long:  `Creates an env/<name>.yaml file describing the environment and target folder (configs/<name>).`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var name string

		// Priority: arg, flag, interactive prompt
		if len(args) > 0 {
			name = args[0]
		}

		if name == "" {
			nameFlag, _ := cmd.Flags().GetString("name")
			if nameFlag != "" {
				name = nameFlag
			}
		}

		if name == "" {
			fmt.Print("Enter environment name: ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("failed to read input: %v", err)
			}
			name = strings.TrimSpace(input)
		}

		if name == "" {
			fmt.Println("Environment name is required")
			return
		}

		filePath, err := env.CreateEnv(name)
		if err != nil {
			log.Fatalf("failed to create env file: %v", err)
		}

		fmt.Printf("Created env file: %s\n", filePath)
	},
}

func init() {
	rootCmd.AddCommand(generateEnvCmd)
	generateEnvCmd.Flags().StringP("name", "n", "", "Name of the environment to create")
}
