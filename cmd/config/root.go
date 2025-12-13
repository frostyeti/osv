/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var ConfigCmd = &cobra.Command{
	Use:   "config ",
	Short: "Manage configuration settings",
	Long: `Set one config value in the config file.

Examples:
  # Sets the default service name
  osv config set service my-service-name
  
  # Gets the default service name
  osv config get service`,
	Run: func(cmd *cobra.Command, args []string) {

		cmd.Help()
	},
}

func init() {
	ConfigCmd.AddCommand(getCmd)
	ConfigCmd.AddCommand(setCmd)
	ConfigCmd.AddCommand(rmCmd)
}
