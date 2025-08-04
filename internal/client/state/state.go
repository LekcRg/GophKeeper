package state

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/LekcRg/GophKeeper/internal/client/req"
	"github.com/LekcRg/GophKeeper/internal/config"
	"github.com/LekcRg/GophKeeper/internal/crypto"
	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	"github.com/charmbracelet/bubbles/table"
)

type State struct {
	req       *req.Request
	config    *config.ClientConfig
	Vault     []models.VaultItem
	Table     []table.Row
	CryptoKey []byte
	activeID  int
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
	for i := range s.Vault {
		item := &s.Vault[i]

		id := strconv.Itoa(item.ID)
		formattedDate := item.UpdatedAt.Format("2 January 2006")
		s.Table[i] = table.Row{id, item.Name, item.Type, formattedDate}
	}
}

func (s *State) SaveCryptoPassword(cryptoKey []byte) {
	s.CryptoKey = cryptoKey
}

func (s *State) AddVaultItem(item models.VaultItem) {
	s.Vault = append(s.Vault, item)

	s.updateTable()
}

func (s *State) SetActiveID(id int) {
	s.activeID = id
}

func (s *State) GetActiveItem() (models.VaultItem, error) {
	for i := range s.Vault {
		item := &s.Vault[i]
		if item.ID == s.activeID {
			enc, err := base64.StdEncoding.DecodeString(item.EncryptedDataString)
			if err != nil {
				return *item, err
			}

			decryptedJSON, err := crypto.Decrypt(enc, s.CryptoKey)
			if err != nil {
				return *item, err
			}

			newItem := *item

			switch item.Type {
			case "password":
				var tmp models.VaultItemDataPassword
				err = json.Unmarshal(decryptedJSON, &tmp)
				newItem.DecryptedData = tmp
			case "note":
				var tmp models.VaultNote
				err = json.Unmarshal(decryptedJSON, &tmp)
				newItem.DecryptedData = tmp
			case "card":
				var tmp models.VaultItemDataCard
				err = json.Unmarshal(decryptedJSON, &tmp)
				newItem.DecryptedData = tmp
			case "binary":
				var tmp models.VaultItemDataBinary
				err = json.Unmarshal(decryptedJSON, &tmp)
				newItem.DecryptedData = tmp
			default:
				return models.VaultItem{}, fmt.Errorf("unknown item type %q", item.Type)
			}

			if err != nil {
				return models.VaultItem{}, err
			}

			return newItem, nil
		}
	}

	return models.VaultItem{}, errs.ErrNotFourndActiveItem
}
