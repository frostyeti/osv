/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package config

import (
	"os"
	"path/filepath"

	cfg "github.com/frostyeti/osv/internal/config"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get one config value from the config file",
	Long: `Get one config value from the config file.

Examples:
  # Gets the default service name
  osv config get service`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			Error(cmd, "key must be provided\n")
			os.Exit(1)
		}

		key := args[0]

		if key == "" {
			Error(cmd, "key must be not be empty\n")
			os.Exit(1)
		}

		configFile, err := GetConfigPath()
		if err != nil {
			Error(cmd, "getting config path: %v\n", err)
			os.Exit(1)
		}

		_, err = os.Stat(configFile)
		if err != nil && os.IsNotExist(err) {
			parent := filepath.Dir(configFile)
			if _, err := os.Stat(parent); os.IsNotExist(err) {
				err = os.MkdirAll(parent, 0755)
				if err != nil {
					Error(cmd, "creating config directory: %v\n", err)
					os.Exit(1)
				}
			}

			err2 := os.WriteFile(configFile, []byte{}, 0644)
			if err2 != nil {
				Error(cmd, "creating config file: %v\n", err2)
				os.Exit(1)
			}
		}

		if err != nil {
			Error(cmd, "stating config file: %v\n", err)
			os.Exit(1)
		}

		kv := cfg.NewConfig()
		err = kv.Load(configFile)
		if err != nil {
			Error(cmd, "loading config file: %v\n", err)
			os.Exit(1)
		}

		value, ok := kv.Get(key)
		if !ok {
			Error(cmd, "%s not set\n", key)
			os.Exit(1)
		}

		os.Stdout.WriteString(value + "\n")
	},
}
