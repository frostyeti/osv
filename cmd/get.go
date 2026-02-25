/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/frostyeti/go/dotenv"
	"github.com/frostyeti/osv/internal/utils"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <key>...",
	Short: "Get one or more secrets from the keyring",
	Long: `Get one or more secrets from the OS keyring.

Examples:
  # Get a single secret
  osv --service <service-name> get --key my-secret

  # or set the service name via environment variable 
  export OSV_SERVICE=<service-name>
  osv get --key my-secret

  # Get multiple secrets
  osv get --key secret1 --key secret2

  # Get secrets with different output formats
  osv get secret1 --format json
  osv get --key secret1 --format sh
  osv get --key secret1 --format dotenv`,

	Run: func(cmd *cobra.Command, args []string) {
		keys, _ := cmd.Flags().GetStringSlice("key")
		format, _ := cmd.Flags().GetString("format")
		clip, _ := cmd.Flags().GetBool("clip")

		if len(args) > 0 {
			keys = append(keys, args...)
		}

		if format == "" {
			format = "text"
		}

		if len(keys) == 0 {
			Error(cmd, "at least one --key must be provided\n")
			osExit(1)
		}

		kr, err := openKeyring(cmd)
		if err != nil {
			Error(cmd, "opening keyring failed: %v\n", err)
			osExit(1)
		}

		values := map[string]string{}
		var firstVal string
		for i, key := range keys {
			item, err := kr.Get(key)
			if err != nil {
				Error(cmd, "getting secret %s failed: %v\n", key, err)
				osExit(1)
			}
			val := string(item.Data)
			values[key] = val
			if i == 0 {
				firstVal = val
			}
		}

		if clip {
			if err := clipboard.WriteAll(firstVal); err != nil {
				Error(cmd, "copying to clipboard failed: %v\n", err)
				osExit(1)
			}
			Ok(cmd, "copied to clipboard\n")
			return
		}

		switch format {
		case "json":
			b, err := json.MarshalIndent(values, "", "  ")
			if err != nil {
				Error(cmd, "marshaling secrets to JSON failed: %v\n", err)
				osExit(1)
			}
			fmt.Println(string(b))

		case "null-terminated", "null":
			for _, v := range values {
				fmt.Printf("%s\x00", v)
			}

		case "sh", "bash", "zsh":
			for k, v := range values {
				key := utils.ScreamingSnakeCase(k)
				fmt.Printf("export %s='%s'\n", key, v)
			}

		case "powershell", "pwsh":
			for k, v := range values {
				key := utils.ScreamingSnakeCase(k)
				fmt.Printf("$Env:%s='%s'\n", key, v)
			}

		case "dotenv", "env", ".env":
			doc := dotenv.NewDoc()
			for k, v := range values {
				key := utils.ScreamingSnakeCase(k)
				doc.Set(key, v)
			}
			fmt.Println(doc.String())

		case "azure-devops", "ado":
			for k, v := range values {
				key := utils.ScreamingSnakeCase(k)
				fmt.Printf("##vso[task.setvariable variable=%s;issecret=true]%s\n", key, v)
			}

		case "github":
			for k, v := range values {
				key := utils.ScreamingSnakeCase(k)
				fmt.Printf("::add-mask::%s\n", v)
				envPath := os.Getenv("GITHUB_ENV")
				if envPath == "" {
					Error(cmd, "GITHUB_ENV environment variable is not set\n")
					osExit(1)
				}
				f, err := os.OpenFile(envPath, os.O_APPEND|os.O_WRONLY, 0600)
				if err != nil {
					Error(cmd, "opening GITHUB_ENV file failed: %v\n", err)
					osExit(1)
				}
				defer f.Close()
				if strings.ContainsAny(v, "\r\n") {
					_, _ = f.WriteString(fmt.Sprintf("%s<<EOF\n%s\nEOF\n", key, v))
				} else {
					_, _ = f.WriteString(fmt.Sprintf("%s=%s\n", key, v))
				}
			}

		default:
			for _, v := range values {
				fmt.Println(v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	service := os.Getenv("OSV_SERVICE")
	getCmd.Flags().StringP("service", "s", service, "Service name for the keyring")
	getCmd.Flags().StringSliceP("key", "k", []string{}, "Name of secret(s) to get (can be specified multiple times)")
	getCmd.Flags().StringP("format", "f", "text", "Output format (text, json, sh, bash, zsh, powershell, pwsh, dotenv)")
	getCmd.Flags().BoolP("clip", "c", false, "Copy the first secret to clipboard instead of printing")
}
