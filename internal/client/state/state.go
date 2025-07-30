package state

import (
	"context"
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/table"
)

type State struct {
	ActiveVaultItem string
	Vault           []models.VaultItem
	req             *req.Request
	config          *config.ClientConfig
	Table           []table.Row
}

func New(r *req.Request, cfg *config.ClientConfig) *State {
	return &State{
		req:    r,
		config: cfg,
	}
}

func (s *State) LoadVault(ctx context.Context) error {
	var err error

	s.Vault, err = s.req.VaultGetAll(ctx)
	if err != nil {
		return err
	}

	s.updateTable()

	return nil
}

func (s *State) updateTable() {
	s.Table = make([]table.Row, len(s.Vault))
	for i, item := range s.Vault {
		id := strconv.Itoa(item.ID)
		formattedDate := item.UpdatedAt.Format("2 January 2006")
		s.Table[i] = table.Row{id, item.Name, item.Type, formattedDate}
	}
}
