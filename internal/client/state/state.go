package state

import (
	"github.com/LekcRg/GophKeeper/internal/client/repository"
	"github.com/LekcRg/GophKeeper/internal/config"
)

type State struct {
	repo            *repository.Repository
	ActiveVaultItem string
	Config          config.ClientConfig
	Vault           []string
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
