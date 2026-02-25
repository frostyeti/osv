/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage local osv configuration settings",
	Long: `Manage local osv configuration settings using the CLI.

Examples:
  # Set the default service name
  osv config set service my-service-name
  
  # Get the default service name
  osv config get service`,
	Run: func(cmd *cobra.Command, args []string) {

		_ = cmd.Help()
	},
}

func init() {
	ConfigCmd.AddCommand(getCmd)
	ConfigCmd.AddCommand(setCmd)
	ConfigCmd.AddCommand(rmCmd)
}
