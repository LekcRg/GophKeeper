package config

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

type ClientConfig struct {
	Address            string `yaml:"address"`
	Key                string `yaml:"key"`
	EncryptedTagString string `yaml:"encrypted_tag"`
	SaltString         string `yaml:"salt"`
	EnctyptedTag       []byte `yaml:"-"`
	Salt               []byte `yaml:"-"`
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
	cfg := &ClientConfig{}

	path, err := getClientPath()
	if err != nil {
		return cfg, err
	}

	cfgPath := filepath.Join(path, configFileName)

	f, err := os.Open(cfgPath)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	err = yaml.NewDecoder(f).Decode(&cfg)
	if err != nil {
		return cfg, err
	}

	if cfg.SaltString != "" {
		cfg.Salt, err = base64.StdEncoding.DecodeString(cfg.SaltString)
		if err != nil {
			return cfg, err
		}
	}

	if cfg.EncryptedTagString != "" {
		cfg.EnctyptedTag, err = base64.StdEncoding.DecodeString(cfg.EncryptedTagString)
		if err != nil {
			return cfg, err
		}
	}

	return cfg, nil
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

	if len(c.Salt) > 0 {
		c.SaltString = base64.StdEncoding.EncodeToString(c.Salt)
	}

	if len(c.EnctyptedTag) > 0 {
		c.EncryptedTagString = base64.StdEncoding.EncodeToString(c.EnctyptedTag)
	}

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
