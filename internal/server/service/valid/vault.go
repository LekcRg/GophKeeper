package valid

import (
	"slices"

	"github.com/LekcRg/GophKeeper/internal/errs"
	"github.com/LekcRg/GophKeeper/internal/models"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func vaultType(value any) error {
	s, _ := value.(string)
	types := []string{"password", "note", "card", "binary"}

	if !slices.Contains(types, s) {
		return errs.ErrVaultNotCorrectType
	}

	return nil
}

func VaultCreateItem(item *models.VaultItem) error {
	return validation.ValidateStruct(item,
		validation.Field(&item.Name, validation.Required),
		validation.Field(&item.Type,
			validation.Required, validation.By(vaultType)),
		validation.Field(&item.EncryptedData, validation.Required),
	)
}
