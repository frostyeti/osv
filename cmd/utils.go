/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/99designs/keyring"
	"github.com/frostyeti/osv/cmd/config"
	"github.com/spf13/cobra"
)

var KeyringProvider = defaultOpenKeyring

func openKeyring(cmd *cobra.Command) (keyring.Keyring, error) {
	return KeyringProvider(cmd)
}

func defaultOpenKeyring(cmd *cobra.Command) (keyring.Keyring, error) {
	service, _ := cmd.Flags().GetString("service")
	cfg, confErr := config.GetConfig()

	if service == "" {

		if confErr == nil {
			v, ok := cfg.Get("service")
			if ok {
				service = v
			}
		}
	}

	libSecret := "login"
	keychain := "login"

	if confErr == nil {
		v, ok := cfg.Get("libsecret.collection")
		if ok {
			libSecret = v
		}

		v, ok = cfg.Get("keychain.name")
		if ok {
			keychain = v
		}
	}

	kr, err := keyring.Open(keyring.Config{
		ServiceName:             service,
		LibSecretCollectionName: libSecret,
		KeychainName:            keychain,
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.WinCredBackend,
			keyring.SecretServiceBackend,
		},
	})
	return kr, err
}

func Error(cmd *cobra.Command, format string, a ...interface{}) {
	os.Stderr.WriteString("\x1b[31m[error]\x1b[0m ")
	cmd.PrintErrf(format, a...)
}

func Warning(cmd *cobra.Command, format string, a ...interface{}) {
	os.Stderr.WriteString("\x1b[33m[warning]\x1b[0m ")
	cmd.PrintErrf(format, a...)
}

func Ok(cmd *cobra.Command, format string, a ...interface{}) {
	os.Stderr.WriteString("\x1b[32m[ok]\x1b[0m ")
	cmd.Printf(format, a...)
}

func toScreamingSnakeCase(input string) string {
	output := ""
	for i, char := range input {
		if char >= 'A' && char <= 'Z' {
			if i > 0 {
				output += "_"
			}
			output += string(char)
		} else if char >= 'a' && char <= 'z' {
			if i > 0 && input[i-1] >= 'A' && input[i-1] <= 'Z' {
				output += "_"
			}
			output += string(char - ('a' - 'A'))
		} else if char >= '0' && char <= '9' {
			output += string(char)
		} else {
			if i > 0 && input[i-1] != '_' {
				output += "_"
			}
		}
	}
	return output
}
