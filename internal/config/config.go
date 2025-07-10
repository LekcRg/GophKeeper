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

type Postgres struct {
	User     string `yaml:"user" env:"POSTGRES_USER" long:"pg-user" description:"Postgress user"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" long:"pg-pass" description:"Postgress password"`
	Host     string `yaml:"host" env:"POSTGRES_HOST" long:"pg-host" description:"Postgress host"`
	Port     string `yaml:"port" env:"POSTGRES_PORT" long:"pg-port" description:"Postgress port"`
	DB       string `yaml:"db" env:"POSTGRES_DB" long:"pg-db" description:"Postgress database name"`
	URI      string `yaml:"uri" env:"POSTGRES_URI" long:"pg-uri" description:"Postgress URI"`
	MaxConns string `yaml:"max_conns" env:"MAX_CONNS" long:"pg-max-conns" description:"Postgress max poll connection"`
}

type Config struct {
	Config string `env:"CONFIG" short:"c" long:"config" description:"Path to yaml config"`
	IsDev  bool   `yaml:"is_dev" env:"IS_DEV" short:"d" long:"dev" description:"Dev mode"`
	Addr   string `yaml:"address" env:"ADDRESS" short:"a" long:"addresss" description:"Address for HTTP server"`

	Postgres Postgres `yaml:"postgres"`
}

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
	return &Config{
		Addr: "localhost:8080",
		Postgres: Postgres{
			Host:     "localhost",
			Port:     "5432",
			MaxConns: "10",
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
