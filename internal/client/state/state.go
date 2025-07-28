package state

import (
	"github.com/LekcRg/GophKeeper/internal/client/repository"
	"github.com/LekcRg/GophKeeper/internal/config"
)

type State struct {
	Config          config.ClientConfig
	Vault           []string
	ActiveVaultItem string
	repo            *repository.Repository
}

func New(repo *repository.Repository) *State {
	return &State{
		repo: repo,
	}
}

func (s *State) UpdateEncryptionData(tag, salt []byte, key string) error {
	s.Config.EnctyptedTag = tag
	s.Config.Salt = salt
	s.Config.Key = key

	// return s.Repo.SaveConfig(s.Config)
	return nil
}
