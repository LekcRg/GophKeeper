package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	Address      string `yaml:"address"`
	Key          string `yaml:"key"`
	EnctyptedTag []byte `yaml:"encrypted_tag"`
	Salt         []byte `yaml:"salt"`
}

const (
	configFolder   = "GophKeeper"
	configFileName = "config.yml"
)

func getMacClientPath() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" && filepath.IsAbs(xdg) {
		return xdg, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config"), nil
}

func getClientPath() (string, error) {
	var (
		base string
		err  error
	)

	if runtime.GOOS == "darwin" {
		base, err = getMacClientPath()
		if err != nil {
			return "", err
		}
	} else {
		dir, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}

		base = dir
	}

	return filepath.Join(base, configFolder), nil
}

func GetClientConfig() (*ClientConfig, error) {
	emptyCfg := &ClientConfig{}

	path, err := getClientPath()
	if err != nil {
		return emptyCfg, err
	}

	cfgPath := filepath.Join(path, configFileName)

	f, err := os.Open(cfgPath)
	if err != nil {
		return emptyCfg, err
	}
	defer f.Close()

	var cfg ClientConfig

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return emptyCfg, err
	}

	return &cfg, nil
}

func (c *ClientConfig) updateConfigFile() error {
	path, err := getClientPath()
	if err != nil {
		return err
	}

	var (
		dirPerm  os.FileMode = 0o700
		filePerm os.FileMode = 0o600
	)

	err = os.MkdirAll(path, dirPerm)
	if err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	configPath := filepath.Join(path, configFileName)

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, filePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	err = yaml.NewEncoder(f).Encode(c)
	if err != nil {
		return fmt.Errorf(
			"failed to save config.\nTry deleting the file at:\n  %s\n"+
				"and restarting the app, or report this issue.\n\nDetails: %w",
			path, err,
		)
	}

	return nil
}

func (c *ClientConfig) Update(f func(cfg *ClientConfig)) error {
	f(c)

	return c.updateConfigFile()
}

func wrapSaveErr(path string, err error) error {
	return fmt.Errorf("failed to save config.\nTry deleting the file at:\n  %s\nand restarting the app, or report this issue.\n\nDetails: %w", path, err)
}
