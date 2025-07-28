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

const configFileName = "config.yml"

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

	return filepath.Join(base, "GophKeeper"), nil
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

	err = os.MkdirAll(path, 0700)
	if err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	configPath := filepath.Join(path, configFileName)

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	return yaml.NewEncoder(f).Encode(c)
}

func (c *ClientConfig) Update(f func(cfg *ClientConfig)) error {
	f(c)

	return c.updateConfigFile()
}
