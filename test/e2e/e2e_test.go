//go:build integration
// +build integration

package e2e

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var cliPath string

func TestMain(m *testing.M) {
	// Build the CLI binary before running tests
	tmpDir, err := os.MkdirTemp("", "osv-e2e-*")
	if err != nil {
		panic("failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	cliPath = filepath.Join(tmpDir, "osv")

	// Build the binary
	cmd := exec.Command("go", "build", "-o", cliPath, "../../main.go")
	if err := cmd.Run(); err != nil {
		panic("failed to build CLI: " + err.Error())
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func runCLI(args ...string) (string, string, error) {
	cmd := exec.Command(cliPath, args...)
	// Use a test service to isolate keyring items
	cmd.Env = append(os.Environ(), "OSV_SERVICE=osv-e2e-test-service")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestE2ESecretsWorkflow(t *testing.T) {
	// Step 1: Ensure clean state for the test key
	runCLI("rm", "e2e-test-secret", "--yes")

	// Step 2: Set a secret
	stdout, stderr, err := runCLI("set", "--key", "e2e-test-secret", "--value", "e2e-super-secret-value")
	// If the OS keyring is not available or locked (like in headless CI without dbus-run-session),
	// we just skip the rest of the test or log a warning rather than failing the build.
	if err != nil {
		t.Skipf("Skipping E2E test due to execution failure (keyring might be locked/unavailable): %v\nStderr: %s", err, stderr)
		return
	}
	if !strings.Contains(stdout, "e2e-test-secret is set") && !strings.Contains(stderr, "is set") {
		t.Errorf("Unexpected set output: stdout=%s stderr=%s", stdout, stderr)
	}

	// Step 3: Get the secret
	stdout, stderr, err = runCLI("get", "e2e-test-secret")
	if err != nil {
		t.Errorf("Failed to get secret: %v, stderr: %s", err, stderr)
	}
	if strings.TrimSpace(stdout) != "e2e-super-secret-value" {
		t.Errorf("Expected 'e2e-super-secret-value', got '%s'", stdout)
	}

	// Step 4: Rename the secret
	stdout, stderr, err = runCLI("rename", "e2e-test-secret", "e2e-test-secret-renamed")
	if err != nil {
		t.Errorf("Failed to rename secret: %v, stderr: %s", err, stderr)
	}

	// Step 5: Verify old secret is gone
	_, _, err = runCLI("get", "e2e-test-secret")
	if err == nil {
		t.Errorf("Expected getting old secret to fail")
	}

	// Step 6: Verify new secret exists
	stdout, stderr, err = runCLI("get", "e2e-test-secret-renamed")
	if err != nil {
		t.Errorf("Failed to get renamed secret: %v, stderr: %s", err, stderr)
	}
	if strings.TrimSpace(stdout) != "e2e-super-secret-value" {
		t.Errorf("Expected renamed secret to match 'e2e-super-secret-value', got '%s'", stdout)
	}

	// Step 7: List secrets
	stdout, stderr, err = runCLI("ls", "e2e-*")
	if err != nil && err.Error() != "exit status 1" { // ls can return 1 if no match, but we expect match
		t.Errorf("Failed to ls secrets: %v, stderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "e2e-test-secret-renamed") {
		t.Errorf("Expected ls to include renamed secret. Got: %s", stdout)
	}

	// Step 8: Remove the secret
	stdout, stderr, err = runCLI("rm", "e2e-test-secret-renamed", "--yes")
	if err != nil {
		t.Errorf("Failed to remove secret: %v, stderr: %s", err, stderr)
	}
	if !strings.Contains(stdout, "Successfully deleted") {
		t.Errorf("Expected success message for rm. Got: %s", stdout)
	}
}
