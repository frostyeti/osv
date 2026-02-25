/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/frostyeti/osv/cmd/config"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "osv",
	Version: Version,
	Short:   "Operating System Vaults (osv) CLI",
	Long: `osv is a secure, cross-platform CLI tool to interact with your operating system's native keychain and credential vaults.

It allows you to securely set, get, list, remove, and rename secrets using your native OS credential store (macOS Keychain, Windows Credential Manager, Linux Secret Service / Keyring).

Examples:
  # Set a secret value
  osv set --key api-token --value "super-secret"

  # Retrieve a secret
  osv get api-token

  # Retrieve multiple secrets and export them as .env content
  osv get api-token db-password --format dotenv

  # List secrets
  osv ls "api-*"`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		osExit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.osv.yaml)")
	rootCmd.AddCommand(config.ConfigCmd)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Print("test")
			return nil
		},
	})
}
