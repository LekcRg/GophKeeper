package req

import (
	"context"
	"errors"
	"fmt"

	"github.com/LekcRg/GophKeeper/internal/models"
)

func (r *Request) CreateVaultItem(
	ctx context.Context, token string, item models.VaultCreateItemReq,
) (models.VaultItem, error) {
	var (
		res     models.VaultItem
		resErrs map[string]string
	)

	_, err := r.client.R().
		SetBody(item).
		SetResult(&res).
		SetError(&resErrs).
		Post("http://localhost:8080/vault/create")
	if err != nil {
		return models.VaultItem{}, err
	}

	if len(resErrs) > 0 {
		fmt.Println(resErrs)
		return models.VaultItem{}, errors.New("req error")
	}

	return res, nil
}
