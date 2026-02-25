package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/99designs/keyring"
	"github.com/spf13/cobra"
)

type mockKeyring struct {
	items map[string]keyring.Item
}

func (m *mockKeyring) Get(key string) (keyring.Item, error) {
	if i, ok := m.items[key]; ok {
		return i, nil
	}
	return keyring.Item{}, keyring.ErrKeyNotFound
}

func (m *mockKeyring) Set(item keyring.Item) error {
	if m.items == nil {
		m.items = make(map[string]keyring.Item)
	}
	m.items[item.Key] = item
	return nil
}

func (m *mockKeyring) Remove(key string) error {
	if _, ok := m.items[key]; !ok {
		return keyring.ErrKeyNotFound
	}
	delete(m.items, key)
	return nil
}

func (m *mockKeyring) Keys() ([]string, error) {
	var keys []string
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys, nil
}

func (m *mockKeyring) GetMetadata(key string) (keyring.Metadata, error) {
	return keyring.Metadata{}, nil
}

func (m *mockKeyring) SetMetadata(key string, metadata keyring.Metadata) error {
	return nil
}

func setupTest(t *testing.T) (*mockKeyring, *bytes.Buffer, *bytes.Buffer) {
	mk := &mockKeyring{items: make(map[string]keyring.Item)}
	KeyringProvider = func(cmd *cobra.Command) (keyring.Keyring, error) {
		return mk, nil
	}

	var exitCode int
	_ = exitCode
	osExit = func(code int) {
		exitCode = code
		panic("osExit called")
	}

	outBuf := new(bytes.Buffer)
	errBuf := new(bytes.Buffer)

	rootCmd.SetOut(outBuf)
	rootCmd.SetErr(errBuf)

	t.Cleanup(func() {
		KeyringProvider = defaultOpenKeyring
		osExit = os.Exit
		rootCmd.SetOut(os.Stdout)
		rootCmd.SetErr(os.Stderr)
	})

	return mk, outBuf, errBuf
}

func executeCommand(args ...string) (string, string, error) {
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
					// normal
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
	_, _ = outBuf.ReadFrom(rOut)
	_, _ = errBuf.ReadFrom(rErr)

	return outBuf.String(), errBuf.String(), err
}

func TestSetAndGetCmd(t *testing.T) {
	mk, _, _ := setupTest(t)

	// Test Set
	out, errOut, _ := executeCommand("set", "--key", "my-secret", "--value", "super-secret-value")
	if !strings.Contains(out, "my-secret is set") {
		t.Errorf("Expected success output, got: %s, err: %s", out, errOut)
	}

	item, err := mk.Get("my-secret")
	if err != nil {
		t.Fatalf("Secret not set in mock keyring")
	}
	if string(item.Data) != "super-secret-value" {
		t.Errorf("Expected secret value to be super-secret-value, got %s", string(item.Data))
	}

	// Test Get
	out, _, _ = executeCommand("get", "my-secret")
	if strings.TrimSpace(out) != "super-secret-value" {
		t.Errorf("Expected get to output 'super-secret-value', got '%s'", out)
	}
}

func TestLsCmd(t *testing.T) {
	mk, _, _ := setupTest(t)
	_ = mk.Set(keyring.Item{Key: "app-db-pass", Data: []byte("pass1")})
	_ = mk.Set(keyring.Item{Key: "app-api-key", Data: []byte("key1")})
	_ = mk.Set(keyring.Item{Key: "other-secret", Data: []byte("sec1")})

	out, _, _ := executeCommand("ls", "app-*")
	if !strings.Contains(out, "app-db-pass") || !strings.Contains(out, "app-api-key") {
		t.Errorf("Expected to list app-* secrets, got: %s", out)
	}
	if strings.Contains(out, "other-secret") {
		t.Errorf("Expected not to list other-secret, got: %s", out)
	}
}

func TestRmCmd(t *testing.T) {
	mk, _, _ := setupTest(t)
	_ = mk.Set(keyring.Item{Key: "to-remove", Data: []byte("val")})

	out, _, _ := executeCommand("rm", "to-remove", "--yes")
	if !strings.Contains(out, "Successfully deleted 1 secret") {
		t.Errorf("Expected success deletion output, got: %s", out)
	}

	_, err := mk.Get("to-remove")
	if err != keyring.ErrKeyNotFound {
		t.Errorf("Expected secret to be removed from mock keyring")
	}
}

func TestRenameCmd(t *testing.T) {
	mk, _, _ := setupTest(t)
	_ = mk.Set(keyring.Item{Key: "old-sec", Data: []byte("val123")})

	out, errOut, _ := executeCommand("rename", "old-sec", "new-sec")
	if !strings.Contains(out, "renamed old-sec to new-sec") {
		t.Errorf("Expected success rename output, got: %s err: %s", out, errOut)
	}

	_, err := mk.Get("old-sec")
	if err != keyring.ErrKeyNotFound {
		t.Errorf("Expected old-sec to be removed")
	}

	item, err := mk.Get("new-sec")
	if err != nil || string(item.Data) != "val123" {
		t.Errorf("Expected new-sec to contain 'val123'")
	}
}
