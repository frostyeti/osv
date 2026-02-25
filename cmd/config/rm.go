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
var rmCmd = &cobra.Command{
	Use:   "rm <key>",
	Short: "Remove one config value from the config file",
	Long: `Remove one config value from the config file.

Examples:
  # Removes the default service name
  osv config rm <key>`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			Error(cmd, "key must be provided\n")
			osExit(1)
		}

		key := args[0]

		if key == "" {
			Error(cmd, "key must be not be empty\n")
			osExit(1)
		}

		configFile, err := GetConfigPath()
		if err != nil {
			Error(cmd, "getting config path: %v\n", err)
			osExit(1)
		}

		_, err = os.Stat(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				parent := filepath.Dir(configFile)
				if _, err := os.Stat(parent); os.IsNotExist(err) {
					err = os.MkdirAll(parent, 0755)
					if err != nil {
						Error(cmd, "creating config directory: %v\n", err)
						osExit(1)
					}
				}

				err2 := os.WriteFile(configFile, []byte{}, 0644)
				if err2 != nil {
					Error(cmd, "creating config file: %v\n", err2)
					osExit(1)
				}
			} else {
				Error(cmd, " config file: %v\n", err)
				osExit(1)
			}
		}

		kv := cfg.NewConfig()
		err = kv.Load(configFile)
		if err != nil {
			Error(cmd, "loading config file: %v\n", err)
			osExit(1)
		}

		kv.Remove(key)
		err = kv.Save()
		if err != nil {
			Error(cmd, "saving config file: %v\n", err)
			osExit(1)
		}

		os.Stderr.WriteString("[ok] config '" + key + "' removed.\n")
		osExit(0)
	},
}
