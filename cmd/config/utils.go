package config

import (
	"os"
	"path/filepath"

	cfg "github.com/frostyeti/osv/internal/config"
	"github.com/spf13/cobra"
)

func GetConfigPath() (string, error) {
	dir := os.Getenv("OSV_CONFIG_DIR")
	if dir != "" {
		return filepath.Join(dir, "osv.kvc"), nil
	}

	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "osv", "osv.kvc"), nil
}

func GetConfig() (*cfg.Config, error) {
	configFile, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(configFile)
	if err != nil && os.IsNotExist(err) {
		return cfg.NewConfig(), nil
	}

	if err != nil {
		return nil, err
	}

	kv := cfg.NewConfig()
	err = kv.Load(configFile)
	if err != nil {
		return nil, err
	}

	return kv, nil
}

func Error(cmd *cobra.Command, format string, a ...interface{}) {
	os.Stderr.WriteString("\x1b[31m[error]\x1b[0m ")
	cmd.PrintErrf(format, a...)
}
