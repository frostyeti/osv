//go:build integration
// +build integration

package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
)

func setupIntegrationTest(t *testing.T) {
	// Override keyring provider with ArrayKeyring
	ak := keyring.NewArrayKeyring(nil)
	KeyringProvider = func(cmd *cobra.Command) (keyring.Keyring, error) {
		return ak, nil
	}

	// Mock osExit
	var exitCode int
	_ = exitCode
	osExit = func(code int) {
		exitCode = code
		panic("osExit called")
	}

	t.Cleanup(func() {
		KeyringProvider = defaultOpenKeyring
		osExit = os.Exit
	})
}

func executeCmd(args ...string) (string, string, error) {
	rootCmd.SetArgs(args)

	oldStdout := os.Stdout
	oldStderr := os.Stderr
	rOut, wOut, _ := os.Pipe()
	rErr, wErr, _ := os.Pipe()
	os.Stdout = wOut
	os.Stderr = wErr

	rootCmd.SetOut(wOut)
	rootCmd.SetErr(wErr)

	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				if r == "osExit called" {
					// normal exit
				} else {
					panic(r)
				}
			}
		}()
		err = rootCmd.Execute()
	}()

	wOut.Close()
	wErr.Close()
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)
	outBuf.ReadFrom(rOut)
	errBuf.ReadFrom(rErr)

	return outBuf.String(), errBuf.String(), err
}

func TestIntegrationEndToEnd(t *testing.T) {
	setupIntegrationTest(t)

	// 1. Set a secret
	out, errOut, _ := executeCmd("set", "--key", "integration-secret", "--value", "secret123")
	if !strings.Contains(out, "integration-secret is set") {
		t.Fatalf("Failed to set secret: %s %s", out, errOut)
	}

	// 2. Get the secret
	out, errOut, _ = executeCmd("get", "integration-secret")
	if strings.TrimSpace(out) != "secret123" {
		t.Fatalf("Failed to get secret: %s %s", out, errOut)
	}

	// 3. List secrets
	out, errOut, _ = executeCmd("ls", "integration-*")
	if !strings.Contains(out, "integration-secret") {
		t.Fatalf("Failed to list secret: %s %s", out, errOut)
	}

	// 4. Remove secret
	out, errOut, _ = executeCmd("rm", "integration-secret", "--yes")
	if !strings.Contains(out, "Successfully deleted 1 secret") {
		t.Fatalf("Failed to remove secret: %s %s", out, errOut)
	}
}
