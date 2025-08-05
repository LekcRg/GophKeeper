package config

import (
	"errors"
	"io"
	"log"
	"os"

	"dario.cat/mergo"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var ErrNothingMerge = errors.New("nothing to merge")

func merge(dst *Config, cfgs ...*Config) error {
	const minLen = 1
	if len(cfgs) < minLen {
		return ErrNothingMerge
	}

	for _, cfg := range cfgs {
		err := mergo.Merge(dst, cfg, mergo.WithOverride)
		if err != nil {
			return err
		}
	}

	return nil
}

func getYamlCfg(path string, cfg *Config) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(fileBytes, &cfg)
}

func getDefaultCfg() *Config {
	const maxBytesKey = 32

	return &Config{
		Addr: "localhost:8080",
		Postgres: Postgres{
			Host:     "localhost",
			Port:     "5432",
			MaxConns: "10",
		},
		Auth: Auth{
			MaxBytesKey: maxBytesKey,
		},
	}
}

func GetConfig(fl []string) (*Config, error) {
	var (
		err     error
		cfg     = getDefaultCfg()
		flagCfg = &Config{}
		envCfg  = &Config{}
		yamlCfg = &Config{}
	)

	_, err = flags.ParseArgs(flagCfg, fl)
	if err != nil {
		return nil, err
	}

	err = godotenv.Load()
	if err != nil {
		var pathErr *os.PathError
		if errors.As(err, &pathErr) && pathErr.Path == ".env" {
			log.Print(err)
		} else {
			return nil, err
		}
	}

	err = cleanenv.ReadEnv(envCfg)
	if err != nil {
		return nil, err
	}

	yamlPath := flagCfg.Config
	if envCfg.Config != "" {
		yamlPath = envCfg.Config
	}

	if yamlPath != "" {
		err = getYamlCfg(yamlPath, yamlCfg)
		if err != nil {
			return nil, err
		}
	}

	err = merge(cfg, yamlCfg, flagCfg, envCfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
