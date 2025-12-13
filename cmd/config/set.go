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
var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set one config value in the config file",
	Long: `Set one config value in the config file.

Examples:
  # Sets the default service name
  osv config set service my-service-name`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 2 {
			Error(cmd, "key and value must be provided\n")
			os.Exit(1)
		}

		key := args[0]
		value := args[1]

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
		if err != nil {
			if os.IsNotExist(err) {
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
			} else {
				Error(cmd, " config file: %v\n", err)
				os.Exit(1)
			}
		}

		kv := cfg.NewConfig()
		err = kv.Load(configFile)
		if err != nil {
			Error(cmd, "loading config file: %v\n", err)
			os.Exit(1)
		}

		kv.Set(key, value)
		err = kv.Save()
		if err != nil {
			Error(cmd, "saving config file: %v\n", err)
			os.Exit(1)
		}

		os.Stderr.WriteString("[ok] config '" + key + "' updated.\n")
		os.Exit(0)
	},
}
