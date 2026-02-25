/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename <old-key> <new-key>",
	Short: "Rename a secret in the keyring",
	Long: `Rename a secret by copying its value to a new key and deleting the old key.

Examples:
  # Rename my-secret to my-new-secret
  osv rename my-secret my-new-secret`,

	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldKey := args[0]
		newKey := args[1]

		kr, err := openKeyring(cmd)
		if err != nil {
			Error(cmd, "opening keyring failed: %v\n", err)
			osExit(1)
		}

		item, err := kr.Get(oldKey)
		if err != nil {
			Error(cmd, "getting old secret %s failed: %v\n", oldKey, err)
			osExit(1)
		}

		err = kr.Set(keyring.Item{
			Key:  newKey,
			Data: item.Data,
		})
		if err != nil {
			Error(cmd, "setting new secret %s failed: %v\n", newKey, err)
			osExit(1)
		}

		err = kr.Remove(oldKey)
		if err != nil {
			Error(cmd, "removing old secret %s failed: %v\n", oldKey, err)
			osExit(1)
		}

		Ok(cmd, "renamed %s to %s\n", oldKey, newKey)
		osExit(0)
	},
}

func init() {
	rootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringP("service", "s", "", "Service name for the keyring")
}
