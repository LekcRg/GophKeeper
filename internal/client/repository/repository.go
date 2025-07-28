package repository

import "github.com/LekcRg/GophKeeper/internal/config"

type Repository interface {
	SaveConfig(cfg *config.ClientConfig) error
	LoadConfig() (*config.ClientConfig, error)
}
